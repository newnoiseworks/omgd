use std::process::{Command, Stdio};

pub fn run(cmd: &str, cmd_desc: &str, cmd_dir: &str, verbose: bool) {
    println!("{} ...", cmd_desc);

    let cmds: Vec<&str> = cmd
        .split(|c| c == ' ')
        .collect();

    let mut command = Command::new(cmds[0]);

    for n in 1..cmds.len() {
        command.arg(cmds[n]);
    }

    if verbose {
        command.stdout(Stdio::piped());
    }
     
    let mut running_command = command.current_dir(cmd_dir)
            .spawn()
            .unwrap();

    running_command.wait().unwrap();

    println!("Success!")
}

