const fs = require('fs')
const rimraf = require('rimraf')
const util = require('util')
const exec = util.promisify(require('child_process').exec);
const AdmZip = require('adm-zip')

const TMP_DIR = process.cwd() + "/.tmp-tpl-fred/"
const LAUNCHER_DIR = process.cwd() + "../launcher"
const WEBSITE_DIR = process.cwd() + "../website"
const GAME_DIR = process.cwd() + "../game"

const LOCAL_GODOT_WIN_BINARY = process.cwd() + "/../Godot/Godot.exe"
const GODOT_WIN_DOWNLOAD_URL = "https://downloads.tuxfamily.org/godotengine/3.1/mono/Godot_v3.1-stable_mono_win64.zip"

const command = process.argv[2]
const environment = process.argv[3]
const version = process.argv[4]

const is_local = environment === "local"

function setup() {
  setupTempWorkingDirectory()
}

function setupTempWorkingDirectory() {
  if (fs.existsSync(TMP_DIR))
    rimraf.sync(TMP_DIR)
 
  fs.mkdirSync(TMP_DIR)

  if (is_local === false) {
    process.chdir(TMP_DIR)
  }
}

async function build() {
  const original_path = process.cwd()

  console.log("building godot project...")
  process.chdir(GAME_DIR)
  await exec(`${LOCAL_GODOT_WIN_BINARY} --export "windows-${environment}" ${TMP_DIR}/game/tpl.exe --no-window`)

  console.log("compressing godot executables...")
  process.chdir(original_path)
  const gameZip = new AdmZip();
  gameZip.addLocalFolder(TMP_DIR + "/game")
  gameZip.writeZip(TMP_DIR + "/tpl-win.zip")

  // TODO: copy .zip to launcher 
  fs.copyFileSync(TMP_DIR + "/tpl-win.zip", LAUNCHER_DIR + "/tpl-win.zip")
  
  // TODO: copy .zip to website repo w/ appropriate version in filename
    // ??? How do I get that version? most recent tag?
}

function deploy() {
  console.log("running deploy...")
}

function buildAndDeploy() {
  console.log("running build and deploy...")
}

// eventually the file should be just the below and reference modules probably
setup()

switch(command) {
  case "build":
    build()
    break
  case "deploy":
    deploy()
    break
  case "build-and-deploy":
    buildAndDeploy()
    break
}