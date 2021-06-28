# omgd
### Open Multiplayer Game Development

#### Making multiplayer game dev easier, including infrastructure deployment

## Vision

`omgd` aims to make it easy to develop multiplayer games. It's not a game engine - `omgd` is an opinionated boilerplate and code generation tool that wraps around existing game, server, and cloud technologies to make multiplayer game development easier, including the creation and deployment of servers on the cloud.

For engines, we're starting with Godot 3 for the game client, Heroic Labs' Nakama for the server framework, and Google Cloud for ~~our favorite big brother~~ a cloud backend.

The long term vision is to add an [agones](https://agones.dev/site/) layer to the stack with automated server builds, to support "instance" style multiplayer development w/n Godot itself (currently you must write your own server code outside the game engine and w/n Nakam's module system).

## Usage
#### Warning: In development! The below are just proposals

- `$ omgd new [project name]` - creates a new project in a folder named after the first argument
- `$ omgd server-start` - starts the local server
- `$ omgd server-stop` - stops the local server
- `$ omgd build-profiles` - builds launch profiles from profiles folder
- `$ omgd build-templates` - builds templates from resources folder
- `$ omgd build-clients` - builds game clients to game/dist folder
- `$ omgd deploy [profile]` - builds and deploys, including infra, to the profile set in the profiles/ folder
- `$ omgd destroy-infra [profile]` - destroys the remote infrastructure

- `$ omgd codegen [type]` - generates project assets, e.g.:
  - `$ omgd codegen channel [dungeon name]` - generates multiplayer "channel" files across game client and server

## Installation

1. Download the latest release from https://github.com/newnoiseworks/omgd/releases/latest
1. Setup the above OS appropriate file as `omgd` within your `$PATH`, e.g. `~/.local/bin`

Currently you must download the repo and build the binary via a rust install and the associated `cargo` tool

`$ cargo build --release && cp targets/release/omgd ~/.local/bin`

### Prereqs

1. git
1. Docker w/ `docker-compose` cli tool available (Docker for desktop works)
1. `gg` -- A separate tool we've built to aid command automation, download the latest release [here](https://github.com/newnoiseworks/gg/releases/latest).
  1.1. Like `omgd` you'll need to place this in your `$PATH` somewhere, we reccommend and may refer from time to time to `~/.local/bin`


### Prereqs (Deployment)
1. `terraform`
1. `butler` (for itch.io deployments)
1. `gcloud` cli tool (for Google Cloud deployemnts)
  1.1. Your google cloud account needs to have billing enabled
  1.2. You must login on your cli via `gcloud auth login`


## Usage


### Create a new project

`$ omgd new [project-name]` -- this will create a new omgd project within the current directory with the name `project-name`


### Generate a realtime multiplayer channel

`$ omgd codegen channel [channel-name]` -- this will generate several files that will help enable realtime multiplayer communication. Adjust the created `resources/[channel-name]-events.yml` file to create individual "events" that, after editing and making a quick `$ omgd build-templates` call, can be published & subscribed to in your game in GDScript via `[ChannelName]Event.[event_name](args)` method calls and connecting to a signal via `[ChannelName]Event.[event_name].connect` calls. You can setup server authoritative responses before passing a published signal to all connected clients by adjusting the `M.match_loop` and `perform_message_validation` methods within the `[channel_name].lua` file.

Note that you'll also need to make sure the Godot event and manager file are added as autoloaded.

Will make a video outlining this procedure, and more precise documentation as well, as things get a bit more "solid."


### Using the resources and profile folders to drive variable templates


## Development

### Prereqs

1. Rust w/ `cargo` cli tool
1. 
