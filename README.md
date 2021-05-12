# omgd
### Open Multiplayer Game Development

#### Making multiplayer game dev easier, including infrastructure deployment

## Vision

`omgd` aims to make it easy to develop multiplayer games. It's not a game engine - `omgd` is an opinionated boilerplate and code generation tool that wraps around existing game and server engines to make multiplayer game development easier, including the creation and deployment of servers on the cloud.

For engines, we're starting with Godot 3 for the game client, Heroic Labs' Nakama for the server framework, and Google Cloud for ~~our favorite big brother~~ a cloud backend.

## Usage
#### Warning: In development! The below are just proposals

`$ omgd new [project name]` - creates a new project in a folder named after the first argument

`$ omgd server-start` - starts the local server
`$ omgd server-stop` - stops the local server
`$ omgd build-templates` - builds templates from resources data
`$ omgd build-clients` - builds game clients to game/dist folder
`$ omgd deploy [profile]` - builds and deploys, including infra, to the profile set in the profiles/ folder
`$ omgd destroy-infra [profile]` - destroys the remote infrastructure

`$ omgd generate [type]` - generates project assets, e.g.:
- `$ omgd generate dungeon [dungeon name]`
- `$ omgd generate dungeon-event [dungeon name] [event name] [event args]`

## Installation

Currently you must download the repo and build the binary via `cargo` (TODO: insert more details here about getting that into ~/.local/bin perhaps, or some other rust-approved method of running this on example projects)

### Prereqs

1. Comfort with the command line and preferabbly a unix or Mac machine, though windows 10 is supported as well.
2. Docker w/ `docker-compose` cli tool available (Docker for desktop works)
3. Gamedev Glue or `gg` - a custom tool we're building [here](https://github.com/newnoiseworks/gg) to assist in automated build and deployment of game clients and servers


### Prereqs (Deployment)
1. `terraform`
2. `butler` (for itch.io deployments)
3. `gcloud` cli tool (for Google Cloud deployemnts)
  3.1. Your google cloud account needs to have billing enabled
  3.2. You must login on your cli via `gcloud auth login`


## Development

### Prereqs

1. Rust w/ `cargo` cli tool
2. 
