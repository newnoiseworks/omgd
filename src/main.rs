use structopt::StructOpt;

mod repo_resource_fetcher;

// omgd new
// omgd generate thing [argsAsSplat]

#[derive(StructOpt)]
struct Cli {
    #[structopt(subcommand)]
    commands: Command
}

#[derive(StructOpt)]
enum Command {
    /// Spawns a new project yall
    New {
        /// The name of the project, will create a folder in this name in the current directory
        name: String,
    },
}

fn main() {
    let args = Cli::from_args();

    match args.commands {
        Command::New { name } => {
            repo_resource_fetcher::get_directory(name)
        }
    }
}

