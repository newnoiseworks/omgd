use std::process::Command;
use std::env;
use std::fs;

// run_cmd_on_dir simply runs a given command on the directory.
//
// Note: Since this splits a string into arguments, certain chain
// commands e.g. `$ echo 'somestring' >> some_file.txt` may not
// work out so well.
//
// * cmd - The command to run
// * cmd_desc - The description of the command
// * cmd_dir - The directory to run the command on
pub fn run_cmd_on_dir(cmd: &str, cmd_desc: &str, cmd_dir: &str) {
    println!("{} ...", cmd_desc);

    let cmds: Vec<&str> = cmd
        .split(|c| c == ' ')
        .collect();

    let mut command = Command::new(cmds[0]);

    for n in 1..cmds.len() {
        command.arg(cmds[n]);
    }
     
    command.current_dir(cmd_dir)
            .status()
            .unwrap();

    println!("Success!")
}

// get_directory_from_repo downloads sub directories from this repo to the target path.
//
// When running in release mode, it will check the live github repo. Locally, it will
// treat the cwd .git repo as the origin. Either way it will perform a sparse checkout
// downloading only the folder requested. Afterwards it will move it to the target path
// and then delete the repo directory.
//
// * `sub_folder` - The sub folder within this repo to download, typically w/n static
// * `target_path` - The final target directory to copy the sub folder to.
pub fn get_directory_from_repo(sub_folder: &str, target_path: &str) {
    let repo;

    if cfg!(debug_assertions) {
        let path = env::current_dir().expect("Couldn't get cwd").to_str().unwrap().to_string();
        repo = format!("{}/.git", &path);
    } else {
        repo = format!("git@github.com:newnoiseworks/omgd.git");
    }

    let home_dir = env::var("HOME").expect("$HOME env var not set");

    let repo_dir = format!("{}/.omgdtmp/repo", home_dir);
    let mk_repo_dir_cmd = format!("mkdir -p {}", repo_dir);

    // This can't be good for windows?
    run_cmd_on_dir(&mk_repo_dir_cmd, "making temp repo as needed", ".");
    run_cmd_on_dir("git init", "setting up git repo in temp folder", &repo_dir);
    run_cmd_on_dir("git config core.sparseCheckout true", "setting folder to sparse checkout", &repo_dir);

    let sparse_checkout_path = format!("{}/.git/info/sparse_checkout", repo_dir);
    fs::write(sparse_checkout_path, sub_folder).expect("Unable to write sparse checkout file");

    let setup_origin = format!("git remote add -f origin {}", repo);
    run_cmd_on_dir(&setup_origin, "setting up sparse-checkout file", &repo_dir);

    let git_pull_cmd;

    if cfg!(debug_assertions) {
        // TODO: If running locally / via cargo run, the below should be automated somehow
        git_pull_cmd = "git pull origin 2-new-project-generator";
    } else {
        git_pull_cmd = "git pull origin 2-new-project-generator";
        // TODO: on final merge of this branch, go w/ below
        // git_pull_cmd = "git pull origin master";
    }
    
    run_cmd_on_dir(&git_pull_cmd, "getting files from repo", &repo_dir);

    let final_move = format!("mv {}/{} {}", repo_dir, sub_folder, target_path);
    run_cmd_on_dir(&final_move, "moving folder", ".");

    let final_cleanup = format!("rm -rf {}", repo_dir);
    run_cmd_on_dir(&final_cleanup, "cleaning up original folder", ".");
}

