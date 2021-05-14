use std::process::Command;

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

pub fn get_directory_from_repo(record_as: String) {
    if cfg!(debug_assertions) {
        // get the directory from the local filesystem
        println!("Recording to directory {} from local filesystem", record_as)
    } else {
        println!("Recording to directory {} from github repo", record_as)
        // get the directory from github.com
    }

    run_cmd_on_dir("gg run", "building clients in game/dist folder...", ".");
}

