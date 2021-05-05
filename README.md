# gg
### GameDev Glue

#### Automated build and deploy of your game's clients, servers, and server infrastructure

## Vision

`gg` is a task runner meant to "glue" together your game servers' infrastructure setup, game client builds, and client / server deployments via a standardized YML configuration file for each "profile" or environment you wish to deploy to.

Useful as a wrapper tool to setup git branch based deployment to new servers for easier testing. You still need to provide your own build, deploy, and infrastrucure processes, but `gg` will run them for you via a simple command structure directed by a `YML` file.

`gg` also allows for templating, so you can share data between your game and server before building - including infrastructure setup e.g. IP addresses of newly created servers, but also for simpler things like inventory, where item prices are only on the server for instance.

See [`omgd`](https://github.com/newnoiseworks/omgd) for a project which uses the `gg` runner base in setting up projects for you with a preconfigured build, deploy, and infra setup process using Godot, Nakama, and Terraform.

## Usage
#### Warning: In development!

Projects for now must be setup as follows, though the goal is to get to a more configurable state

```
/
- /game         (your game client)
- /server       (your server program)
  - /infra
  - /deploy
- /profiles     ()
- /resources
```
