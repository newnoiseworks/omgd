const fs = require('fs')
const rimraf = require('rimraf')
const util = require('util')
const exec = util.promisify(require('child_process').exec);
const AdmZip = require('adm-zip')

const COMMAND = process.argv[2]
const ENVIRONMENT = process.argv[3]
const VERSION = process.argv[4]

const IS_LOCAL = ENVIRONMENT === "local"

const LAUNCHER_DIR = process.cwd() + "/../launcher"
const WEBSITE_DIR = process.cwd() + "/../website"
const GAME_DIR = process.cwd() + "/../game"

const TMP_DIR = process.cwd() + "/.tmp-tpl-fred/"
const TMP_LAUNCHER_DIR = TMP_DIR + "launcher"
const TMP_WEBSITE_DIR = TMP_DIR + "website"
const TMP_GAME_DIR = TMP_DIR + "game"

const LAUNCHER_REPO = "git@github.com:newnoiseworks/tpl-launcher.git"
const WEBSITE_REPO = "git@github.com:newnoiseworks/tpl-website.git"
const GAME_REPO = "git@github.com:newnoiseworks/not-stardew.git"

const LOCAL_GODOT_WIN_BINARY = process.cwd() + "/../Godot/Godot.exe"
const GODOT_WIN_DOWNLOAD_URL = "https://downloads.tuxfamily.org/godotengine/3.1/mono/Godot_v3.1-stable_mono_win64.zip"

// TODO: The below should be defined via tag or project's package.json
const LAUNCHER_VERSION = "0.0.3"

async function setup() {
  console.log("setting up...")
  await setupTempWorkingDirectory()
}

async function setupTempWorkingDirectory() {
  if (fs.existsSync(TMP_DIR))
    rimraf.sync(TMP_DIR)
 
  fs.mkdirSync(TMP_DIR)

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
  },/*{
    originalRepoDir: GAME_DIR,
    originalRepo: GAME_REPO,
    tmpRepoDir: TMP_GAME_DIR,
  }*/].map((proj) => {
    let launcherVersion
    process.chdir(proj.originalRepoDir)
    
    return exec(`git rev-parse HEAD`)
    .then(({ stderr, stdout }) => {
      launcherVersion = stdout.trim()
    }).then(() => {
      process.chdir(TMP_DIR)
      return exec(`git clone --depth 1 ${proj.originalRepo} ${proj.tmpRepoDir}`)
    }).then(() => {
      if (ENVIRONMENT !== "production") {
        process.chdir(proj.tmpRepoDir)
        return exec(`git reset --hard ${launcherVersion}`)
      }
    })
  }))
}

async function build() {
  const original_path = process.cwd()
  const gameDir = GAME_DIR
  const websiteDir = IS_LOCAL ? WEBSITE_DIR : TMP_WEBSITE_DIR
  const launcherDir = IS_LOCAL ? LAUNCHER_DIR : TMP_LAUNCHER_DIR

  console.log("Building godot game client...")

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
  
  let versionString
  if (!!VERSION)
    versionString = VERSION
  else {
    process.chdir(gameDir)
    await exec("git fetch --tags")
    // ??? For non production builds, do we need to adjust build #? Or should we only worry about that for prod pushes? We've gotta remember to update all the /game files manually as well. This really only matters for testing & for triggering updates w/n the launcher so probably not necessary, the only requirement is that all  version #'s are in sync. may be worth automating the references w/n Godot and Nakama somehow as well.
    versionString = (await exec("git describe --tags --abbrev=0")).stdout.trim()
    process.chdir(original_path)
  }

  console.log("copying godot zip to website...")
  fs.copyFileSync(TMP_DIR + "tpl-win.zip", websiteDir + `/public/static/ThePromisedLand-${versionString}.win.zip`)

  if (IS_LOCAL === false) {
    console.log("packaging launcher...")
    process.chdir(launcherDir)
    await exec(`yarn && yarn package`)

    // TODO: Need to get tagged launcher version
    console.log("copying launcher to website...")
    fs.copyFileSync(
      `release\\The Promised Land - Game Launcher Setup ${LAUNCHER_VERSION}.exe`,
      `${websiteDir}/public/static/ThePromisedLand-Launcher-Setup-${LAUNCHER_VERSION}.${versionString}.exe`
    )
  }

  console.log("building website for firebase...")
  process.chdir(websiteDir)

  await exec("npm install --no-progress")
  await exec("npm run build")
}

async function deploy() {
  console.log("running deploy...")

  if (IS_LOCAL) return

  process.chdir(TMP_WEBSITE_DIR)
  await exec(`firebase deploy`)
}

async function buildAndDeploy() {
  console.log("running build and deploy...")
  await build()
  await deploy()
}

// eventually the file should be just the below and reference modules probably
(async function() {
  await setup()
  
  switch(COMMAND) {
    case "build":
      await build()
      break
    case "deploy":
      deploy()
      break
    case "build-and-deploy":
      buildAndDeploy()
      break
  }
})()