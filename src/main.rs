use structopt::StructOpt;

mod run_cmd_on_dir;
//mod repo_resource_fetcher;

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
            let cmd = format!("git clone git@github.com:newnoiseworks/tpl-game-gd.git {}", name);
            let desc = format!("generating new project in folder ./{}", name);
            run_cmd_on_dir::run(&cmd, &desc, ".");
        }
        Command::Codegen { plan } => {
            println!("We will generate code from plan {}.", plan)
        }
        Command::ServerStart { tail } => {
            if tail {
                run_cmd_on_dir::run("docker-compose up", "starting omgd servers...", "server");
            } else {
                run_cmd_on_dir::run("docker-compose up -d", "starting omgd servers...", "server");
            }
        }
        Command::ServerStop {} => {
            run_cmd_on_dir::run("docker-compose down", "stopping omgd servers...", "server");
        }
        Command::BuildTemplates { profile } => {
            match profile {
                Some(p) => {
                    let cmd = format!("gg build-templates . --profile=profiles/{}", p);
                    run_cmd_on_dir::run(&cmd, "building templates...", ".");
                }
                None => run_cmd_on_dir::run("gg build-templates .", "building templates...", ".")
            }
        }
        Command::BuildClients { } => {
            run_cmd_on_dir::run("gg run", "building clients in game/dist folder...", ".");
        }
        Command::Deploy { profile } => {
            run_cmd_on_dir::run("mkdir .omgdtmp", "creating temporary dir...", ".");

            let dir = format!(".omgdtmp/{}", profile);
            let git_clone_cmd = format!("git clone . {}", dir);

            run_cmd_on_dir::run(&git_clone_cmd, "cloning repo...", ".");

            let bnd_cmd = format!("gg run --profile=profiles/{}", profile);
            run_cmd_on_dir::run(&bnd_cmd, "build and deploying repo...", &dir);
        }
        Command::DestroyInfra { profile } => {
            let dir = format!(".omgdtmp/{}", profile);
            let cmd = format!("gg run task destroy-infra --profile=profiles/{}", profile);
            run_cmd_on_dir::run(&cmd, "destroying infra...", &dir);
        }
    }
}

