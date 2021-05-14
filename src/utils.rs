use std::process::Command;
use std::env;
use std::fs;

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

pub fn get_directory_from_repo(sub_folder: &str, target_path: &str) {
    let mut repo = "git@github.com:newnoiseworks/omgd.git";

    if cfg!(debug_assertions) {
        // repo = "."; 
    }

    let home_dir;

    match env::var("HOME") {
        Ok(h) => home_dir = h,
        Err(e) => panic!("No $HOME var set! {}", e),
    }

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
    
    run_cmd_on_dir("git pull origin 2-new-project-generator", "getting files from repo", &repo_dir);

    let final_move = format!("mv {}/{} {}", repo_dir, sub_folder, target_path);

    run_cmd_on_dir(&final_move, "moving folder", ".");

    let final_cleanup = format!("rm -rf {}", repo_dir);

    run_cmd_on_dir(&final_cleanup, "cleaning up original folder", ".");
}

