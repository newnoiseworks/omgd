use structopt::StructOpt;

mod utils;

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
    /// Stop the server
    ServerStop {},
    /// Build templates
    BuildTemplates {
        /// optionally reference a profile by it's name in the profiles directory 
        /// e.g. "staging" == "profiles/staging.yml"
        profile: Option<String>
    },
    /// Build clients -- local only, check game/dist/ after running
    BuildClients {
    },
    /// Deploy the project
    Deploy {
        /// required - reference a profile by it's name in the profiles directory 
        /// e.g. "staging" == "profiles/staging.yml"
        profile: String
    },
    /// Destroys infra for the profile
    DestroyInfra {
        /// required - reference a profile by it's name in the profiles directory 
        /// e.g. "staging" == "profiles/staging.yml"
        profile: String
    },
}

fn main() {
    let args = Cli::from_args();

    match args.commands {
        Command::New { name } => {
            // repo_resource_fetcher::get_directory(name)
            utils::get_directory_from_repo("static/new", &name);
        }
        Command::Codegen { plan } => {
            println!("We will generate code from plan {}.", plan)
        }
        Command::ServerStart { tail } => {
            if tail {
                utils::run_cmd_on_dir("docker-compose up", "starting omgd servers...", "server");
            } else {
                utils::run_cmd_on_dir("docker-compose up -d", "starting omgd servers...", "server");
            }
        }
        Command::ServerStop {} => {
            utils::run_cmd_on_dir("docker-compose down", "stopping omgd servers...", "server");
        }
        Command::BuildTemplates { profile } => {
            match profile {
                Some(p) => {
                    let cmd = format!("gg build-templates . --profile=profiles/{}", p);
                    utils::run_cmd_on_dir(&cmd, "building templates...", ".");
                }
                None => utils::run_cmd_on_dir("gg build-templates .", "building templates...", ".")
            }
        }
        Command::BuildClients { } => {
            utils::run_cmd_on_dir("gg run", "building clients in game/dist folder...", ".");
        }
        Command::Deploy { profile } => {
            utils::run_cmd_on_dir("mkdir .omgdtmp", "creating temporary dir...", ".");

            let dir = format!(".omgdtmp/{}", profile);
            let git_clone_cmd = format!("git clone . {}", dir);

            utils::run_cmd_on_dir(&git_clone_cmd, "cloning repo...", ".");

            let bnd_cmd = format!("gg run --profile=profiles/{}", profile);
            utils::run_cmd_on_dir(&bnd_cmd, "build and deploying repo...", &dir);
        }
        Command::DestroyInfra { profile } => {
            let dir = format!(".omgdtmp/{}", profile);
            let cmd = format!("gg run task destroy-infra --profile=profiles/{}", profile);
            utils::run_cmd_on_dir(&cmd, "destroying infra...", &dir);
        }
    }
}

