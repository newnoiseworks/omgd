# OMGD
### Open Multiplayer Game Development

## Vision

`omgd` aims to make it easy to get into and develop multiplayer games. It's not a game engine - `omgd` is an opinionated boilerplate and code generation tool that wraps around game engines to make multiplayer game development easier, including the creation and deployment of servers on the cloud.

For engines, we're starting with Godot 3 for the game client and Heroic Labs' Nakama for the server framework, and Google Cloud for a server host.


## Usage
#### Warning: In development!


## Installation

### Prereqs

0. Comfort with the command line and preferabbly a unix or Mac machine, though windows 10 is supported as well.
1. Docker w/ `docker-compose` cli tool available
  - Docker for desktop is fine
2. Gamedev Glue or `gg` - a custom tool we're building [here](https://github.com/newnoiseworks/gg) to assist in automated build and deployment of game clients and servers


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
