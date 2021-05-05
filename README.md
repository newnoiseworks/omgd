# gg
### GameDev Glue

#### Automated build and deploy of your game's clients, servers, and server infrastructure

## Vision

`gg` is a task runner meant to "glue" together your game servers' infrastructure setup, game client builds, and client / server deployments via a standardized YML configuration file for each "profile" or environment you wish to deploy to.

Useful as a wrapper tool to setup git branch based deployment to new servers for easier testing. You still need to provide your own build, deploy, and infrastrucure processes, but `gg` will run them for you via a simple command structure directed by a `YML` file.

See [`omgd`](https://github.com/newnoiseworks/omgd) for a project which uses the `gg` runner base in setting up projects for you with a preconfigured build, deploy, and infra setup process using Godot, Nakama, and Terraform.

## Usage
#### Warning: In development!
