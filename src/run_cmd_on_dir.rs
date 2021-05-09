use std::process::{Command, Stdio};

pub fn run(cmd: &str, cmd_desc: &str, cmd_dir: &str, verbose: bool) {
    println!("{} ...", cmd_desc);

    let mut command = Command::new(cmd);

    if verbose {
        command.stdout(Stdio::null());
    }

    command.current_dir(cmd_dir)
            .spawn()
            .expect("Command failed");
}
