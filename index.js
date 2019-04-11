const fs = require('fs')
const rimraf = require('rimraf')
const util = require('util')
const exec = util.promisify(require('child_process').exec);
const AdmZip = require('adm-zip')
const yjs = require('js-yaml')
const md5 = require('md5')

const COMMAND = process.argv[2]
const ENVIRONMENT = process.argv[3]

const IS_LOCAL = ENVIRONMENT === "local"

const LAUNCHER_DIR = process.cwd() + "/../launcher"
const WEBSITE_DIR = process.cwd() + "/../website"
const GAME_DIR = process.cwd() + "/../game"
const SERVER_DIR = process.cwd() + "/../server"

const TMP_DIR = process.cwd() + "/.tmp-tpl-fred/"
const TMP_LAUNCHER_DIR = TMP_DIR + "launcher"
const TMP_WEBSITE_DIR = TMP_DIR + "website"
const TMP_GAME_DIR = TMP_DIR + "game"
const TMP_SERVER_DIR = TMP_DIR + "server"

const LAUNCHER_REPO = "git@github.com:newnoiseworks/tpl-launcher.git"
const WEBSITE_REPO = "git@github.com:newnoiseworks/tpl-website.git"
const GAME_REPO = "git@github.com:newnoiseworks/not-stardew.git"
const SERVER_REPO = "git@github.com:newnoiseworks/not-stardew-backend.git"

const LOCAL_GODOT_WIN_BINARY = process.cwd() + "/../Godot/Godot.exe"
const GODOT_WIN_DOWNLOAD_URL = "https://downloads.tuxfamily.org/godotengine/3.1/mono/Godot_v3.1-stable_mono_win64.zip"

function tresToJsonViaYml(path) {
  return yjs.safeLoad(
    fs.readFileSync(path, 'utf8')
      .replace(/=/g, ':')
      .replace(/ :/g, ':')
      .replace(/\n\n/g, '\n')
      .replace(/\n/g, '\n  ')
      .replace(/  \[/g, '[')
      .replace(/\[/g, '')
      .replace(/\]/g, ':')
  )
}

async function setup() {
  console.log("setting up...")
  await setupTempWorkingDirectory()
}

async function setupTempWorkingDirectory() {
  if (fs.existsSync(TMP_DIR))
    rimraf.sync(TMP_DIR)
 
  fs.mkdirSync(TMP_DIR)
  fs.mkdirSync(TMP_DIR + "/gameBuild")

  if (IS_LOCAL) return

  process.chdir(TMP_DIR)

  console.log("cloning repos...")
  await cloneRepositoriesFromGithub()
}

async function cloneRepositoriesFromGithub() {
  // TODO: abort if "git status" shows changes
  await Promise.all([{
    originalRepoDir: LAUNCHER_DIR,
    originalRepo: LAUNCHER_REPO,
    tmpRepoDir: TMP_LAUNCHER_DIR,
  },{
    originalRepoDir: WEBSITE_DIR,
    originalRepo: WEBSITE_REPO,
    tmpRepoDir: TMP_WEBSITE_DIR,
  },{
    originalRepoDir: SERVER_DIR,
    originalRepo: SERVER_REPO,
    tmpRepoDir: TMP_SERVER_DIR,
  },/*{
    originalRepoDir: GAME_DIR,
    originalRepo: GAME_REPO,
    tmpRepoDir: TMP_GAME_DIR,
  }*/].map((proj) => {
    let originalRepoCommit
    process.chdir(proj.originalRepoDir)
    
    return exec(`git rev-parse HEAD`)
    .then(({ stderr, stdout }) => {
      originalRepoCommit = stdout.trim()
    }).then(() => {
      process.chdir(TMP_DIR)
      return exec(`git clone --depth 5 ${proj.originalRepo} ${proj.tmpRepoDir}`)
    }).then(() => {
      if (ENVIRONMENT !== "production") {
        process.chdir(proj.tmpRepoDir)
        return exec(`git reset --hard ${originalRepoCommit}`)
      }
    })
  }))
}

async function build() {
  const original_path = process.cwd()
  const gameDir = GAME_DIR
  const websiteDir = IS_LOCAL ? WEBSITE_DIR : TMP_WEBSITE_DIR
  const launcherDir = IS_LOCAL ? LAUNCHER_DIR : TMP_LAUNCHER_DIR
  const serverDir = IS_LOCAL ? SERVER_DIR : TMP_SERVER_DIR

  const GAME_VERSION = tresToJsonViaYml(`${GAME_DIR}/Resources/Config/config.tpl.tres`).config.version;
  
  if (IS_LOCAL === false) {
    const filePath = `${GAME_DIR}/Resources/Config/config.tpl_${ENVIRONMENT}.tres`
    const envConfigFile = fs.readFileSync(filePath)
    fs.writeFileSync(filePath, envConfigFile.toString().replace(
      /key \= .+/,
      `key = "${md5(`the-promised-land-${ENVIRONMENT}-v${GAME_VERSION}`)}"`
    ))
  }

  console.log("building godot game client...")

  process.chdir(gameDir)

  if (IS_LOCAL === false) {
    // check git status in staging, git status -s should return empty
    const statusString = (await exec(`git status -s`)).stdout.trim()

    if (!!statusString) throw new Error(`Error! can't build to ${ENVIRONMENT} with unchecked file changes in the game folder. Try git stashing first.`)

    // TODO: in production check current git branch && commit in dir -- should be master && latest
  }

  try {
    await exec(`${LOCAL_GODOT_WIN_BINARY} --export "windows-${ENVIRONMENT}" ${TMP_DIR}gameBuild/tpl.exe --no-window`)
  } catch(err) {
    console.log("!!!")
    console.log(err)
    throw new Error(err)
  }

  console.log("compressing godot executables...")
  process.chdir(original_path)
  const gameZip = new AdmZip();
  gameZip.addLocalFolder(TMP_DIR + "gameBuild")
  gameZip.writeZip(TMP_DIR + "tpl-win.zip")

  console.log("copying godot zip to launcher...")
  fs.copyFileSync(TMP_DIR + "tpl-win.zip", launcherDir + "/tpl-win.zip")

  const gameEnvConfig = tresToJsonViaYml(`${GAME_DIR}/Resources/Config/config.tpl_${ENVIRONMENT}.tres`);

  const LAUNCHER_VERSION = require(`${launcherDir}/package.json`).version

  console.log("copying godot zip to website...")
  fs.copyFileSync(TMP_DIR + "tpl-win.zip", websiteDir + `/public/static/ThePromisedLand-${GAME_VERSION}.win.zip`)

  console.log("creating build-config.json file for all projects...")

  const versionObj = {
    "gameVersion": GAME_VERSION,
    "launcherVersion": LAUNCHER_VERSION,
    "environment": ENVIRONMENT,
    nakama: gameEnvConfig.nakama,
    website: gameEnvConfig.website
  }

  fs.writeFileSync(`${launcherDir}/build-config.json`, JSON.stringify(versionObj))
  fs.copyFileSync(`${launcherDir}/build-config.json`, `${websiteDir}/src/build-config.json`)
  fs.copyFileSync(`${launcherDir}/build-config.json`, `${serverDir}/build-config.json`)

  console.log("packaging launcher...")
  process.chdir(launcherDir)

  const launcherPackageConfig = require(`${launcherDir}/package.json`)
  let productName = launcherPackageConfig.build.productName

  if (IS_LOCAL === false) {
    let appId = launcherPackageConfig.build.appId
    
    if (ENVIRONMENT === "production") {
      productName = productName.replace(" {ENVIRONMENT}", "")
      appId = appId.replace("{ENVIRONMENT}", "")
    } else {
      productName = productName.replace("{ENVIRONMENT}", ENVIRONMENT)
      appId = appId.replace("{ENVIRONMENT}", "-" + ENVIRONMENT)
    }
    
    launcherPackageConfig.build.productName = productName
    launcherPackageConfig.build.appId = appId
    launcherPackageConfig.productName = productName
    fs.writeFileSync(`${launcherDir}/package.json`, JSON.stringify(launcherPackageConfig, null, 2));
  }

  await exec(`yarn && yarn package`)

  // TODO: Need to get tagged launcher version
  console.log("copying launcher to website...")
  fs.copyFileSync(
    `release\\${productName} Setup ${LAUNCHER_VERSION}.exe`,
    `${websiteDir}/public/static/ThePromisedLand-Launcher-Setup-${LAUNCHER_VERSION}.${GAME_VERSION}.${ENVIRONMENT}.exe`
  )

  console.log("building website for firebase...")
  process.chdir(websiteDir)

  await exec("npm install --no-progress")
  await exec("npm run build")

  console.log("building server app with latest versions...")
  process.chdir(serverDir)

  await exec("npm install --no-progress")
  await exec("npm run build")
  
  if (IS_LOCAL){
    console.log("restarting local docker containers")
    await exec("docker-compose down")
    await exec("docker-compose up -d")
  }
}

async function deploy() {
  if (IS_LOCAL) return

  console.log("preparing deploy...")
  console.log("shutting down nakama servers...")
  await exec(`gcloud compute ssh instance-1 --zone us-east1-b --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD:$PWD" -w="$PWD" docker/compose:1.13.0 down"`)

  console.log("running deploy...")

  console.log("deploying website...")
  process.chdir(TMP_WEBSITE_DIR)
  await exec(`firebase deploy`)

  console.log("pushing nakama changes...")
  process.chdir(TMP_SERVER_DIR)
  // push lib and docker compose and build file up to server
  await exec(`gcloud compute scp --zone us-east1-b --recurse --force-key-file-overwrite ./nakama instance-1:`)
  await exec(`gcloud compute scp --zone us-east1-b --force-key-file-overwrite docker-compose.yml instance-1:`)

  // up containers via docker-compose
  await exec(`gcloud compute ssh instance-1 --zone us-east1-b --command "docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v "$PWD:$PWD" -w="$PWD" docker/compose:1.13.0 up -d"`)
}


async function buildAndDeploy() {
  console.log("running build and deploy...")
  await build()
  await deploy()
}

(async function() {
  
  switch(COMMAND) {
    case "build":
      await setup()
      await build()
      break
    case "deploy":
      await deploy()
      break
    case "build-and-deploy":
      await setup()
      await buildAndDeploy()
      break
  }
})()