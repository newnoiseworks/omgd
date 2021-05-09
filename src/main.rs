use structopt::StructOpt;

mod run_cmd_on_dir;
mod repo_resource_fetcher;

#[derive(StructOpt)]
struct Cli {
    #[structopt(subcommand)]
    commands: Command
}

#[derive(StructOpt)]
enum Command {
    /// Spawns a new project
    New {
        /// The name of the project, will create a folder in this name in the current directory
        name: String,
    },
    /// Generate various code
    Codegen {
        /// What you want to generate
        plan: String,
    },
    /// Start the server
    ServerStart {
        /// Set to tail logs, will print output from server to this terminal
        #[structopt(short = "t", long = "tail")]
        tail: bool,
    },
}

fn main() {
    let args = Cli::from_args();

    match args.commands {
        Command::New { name } => {
            repo_resource_fetcher::get_directory(name)
        }
        Command::Codegen { plan } => {
            println!("We will generate code from plan {}.", plan)
        }
        Command::ServerStart { tail } => {
            println!("Start that server");
            if tail {
                run_cmd_on_dir::run("docker-compose up", "start them servers", "server", true);
            } else {
                run_cmd_on_dir::run("docker-compose up -d", "start them servers", "server", false);
            }
        }
    }
}

