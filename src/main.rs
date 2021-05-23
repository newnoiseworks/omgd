use structopt::StructOpt;
use convert_case::{Case, Casing};

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
        args: String,
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
        Command::Codegen { plan, args } => {
            run_codegen_cmds(&plan, &args);
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

pub fn run_codegen_cmds(plan: &str, args: &str) {
    if plan == "channel" {
        run_channel_manager_cmds(args);
    } else if plan == "example-2d-player-movement" {
        run_example_2d_player_movement_cmds(args);
    } else {
        println!("Plan {} not found in codegen commands", plan);
    }
}

fn run_example_2d_player_movement_cmds(channel_type_name: &str) {
    utils::run_cmd_on_dir("mkdir -p .omgdtmp", "making target path if needed", ".");

    let tmp_dir = ".omgdtmp/example-2d-player-movement";

    utils::get_directory_from_repo("static/example-2d-player-movement", &tmp_dir);

    let update_profile = format!("gg update-profile omgd.channel_name {}", &channel_type_name);
    utils::run_cmd_on_dir(&update_profile, "updates temporary profile w/ channel name", &tmp_dir);

    utils::run_cmd_on_dir("gg build-templates . --ext=omgdtpl", "updates temporary profile w/ channel name", &tmp_dir);

    utils::run_cmd_on_dir("rm -rf .omgdtmp/example-2d-player-movement/profiles", "cleaning...", ".");

    let final_move = format!("rsync -a --exclude **/*.omgdtpl {}/ .", &tmp_dir);
    utils::run_cmd_on_dir(&final_move, "rsync'ing folder", ".");

    let final_cleanup = format!("rm -rf {}", &tmp_dir);
    utils::run_cmd_on_dir(&final_cleanup, "cleaning up original folder", ".");

    utils::run_cmd_on_dir("omgd build-templates", "rebuilding .tmpl templates...", ".");
}

fn run_channel_manager_cmds(channel_type_name: &str) {
    utils::run_cmd_on_dir("mkdir -p .omgdtmp", "making target path if needed", ".");

    utils::get_directory_from_repo("static/channel", ".omgdtmp/channel");

    let update_profile = format!("gg update-profile omgd.channel_name {}", channel_type_name);
    utils::run_cmd_on_dir(&update_profile, "updates temporary profile w/ channel name", ".omgdtmp/channel");

    utils::run_cmd_on_dir("gg build-templates . --ext=omgdtpl", "updates temporary profile w/ channel name", ".omgdtmp/channel");

    utils::run_cmd_on_dir("rm -rf .omgdtmp/channel/profiles", "cleaning...", ".");

    let mut change_game_files = format!("mv game/Autoloads/ChannelManager.gd game/Autoloads/{}Manager.gd", channel_type_name.to_case(Case::UpperCamel));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting filenames...", ".omgdtmp/channel");

    change_game_files = format!("sed -i -e s/channel/{}/g game/Autoloads/ChannelEvent.gd.tmpl", channel_type_name.to_case(Case::Snake));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting file contents...", ".omgdtmp/channel");

    change_game_files = format!("sed -i -e s/Channel/{}/g game/Autoloads/ChannelEvent.gd.tmpl", channel_type_name.to_case(Case::UpperCamel));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting file contents...", ".omgdtmp/channel");

    change_game_files = format!("mv game/Autoloads/ChannelEvent.gd.tmpl game/Autoloads/{}Event.gd.tmpl", channel_type_name.to_case(Case::UpperCamel));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting filenames...", ".omgdtmp/channel");

    change_game_files = format!("mv resources/channel_events.yml resources/{}_events.yml", channel_type_name.to_case(Case::Snake));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting filenames...", ".omgdtmp/channel");

    change_game_files = format!("mv game/RootScenes/ChannelMUD.tscn game/RootScenes/{}MUD.tscn", channel_type_name.to_case(Case::UpperCamel));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting filenames...", ".omgdtmp/channel");

    change_game_files = format!("mv game/RootScenes/ChannelMUDController.gd game/RootScenes/{}MUDController.gd", channel_type_name.to_case(Case::UpperCamel));
    utils::run_cmd_on_dir(&change_game_files, "Adjusting filenames...", ".omgdtmp/channel");

    let mut change_server_files = format!("mv server/nakama/data/modules/channel.lua server/nakama/data/modules/{}.lua", channel_type_name.to_case(Case::Snake));
    utils::run_cmd_on_dir(&change_server_files , "Adjusting filenames...", ".omgdtmp/channel");

    change_server_files = format!("mv server/nakama/data/modules/channel_manager.lua server/nakama/data/modules/{}_manager.lua", channel_type_name.to_case(Case::Snake));
    utils::run_cmd_on_dir(&change_server_files , "Adjusting filenames...", ".omgdtmp/channel");
    
    let final_move = format!("rsync -a --exclude **/*.omgdtpl {}/ .", ".omgdtmp/channel");
    utils::run_cmd_on_dir(&final_move, "rsync'ing folder", ".");

    let final_cleanup = format!("rm -rf {}", ".omgdtmp/channel");
    utils::run_cmd_on_dir(&final_cleanup, "cleaning up original folder", ".");

    utils::run_cmd_on_dir("omgd build-templates", "rebuilding .tmpl templates...", ".");
}

