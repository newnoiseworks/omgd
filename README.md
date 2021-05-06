# gg
### GameDev Glue

#### Automated build and deploy of your game's clients, servers, and server infrastructure

## Vision

`gg` is a task runner meant to "glue" together your game servers' infrastructure setup, game client builds, and client / server deployments via a standardized YML configuration file for each "profile" or environment you wish to deploy to, and omitting any of the things you don't need.

Useful as a wrapper tool to setup git branch based deployment to new servers for easier testing. You still need to provide your own build, deploy, and infrastrucure processes, but `gg` will run them for you via a simple command structure directed by a `YML` file.

`gg` also adds resource backed templating, so you can share data between your game and server before building - including infrastructure setup e.g. IP addresses of newly created servers, but also for simpler things like inventory, where item prices are only on the server for instance. Helps to deploy as well as bootstrap local infrastructure for local development.

See [`omgd`](https://github.com/newnoiseworks/omgd) for a project which uses the `gg` runner base in setting up projects for you with a preconfigured build, deploy, and infra setup process using Godot, Nakama, and Terraform.

## Installation

Currently you must download this repo and then build the binary, placing it in your path. For example, given `~/.local/bin` is in your `$PATH` variable, then the following would be good:

`$ go build -o ~/.local/bin/gg`

Eventually we'll post releases with version tags on this repo, and then after that via popular package managers like snap, brew etc.


## Usage
#### Warning: In development!

Projects for now must be setup as follows, though the goal is to get to a more configurable state w/ multiple client, servers, multiple repos, and more.

```
/               (your project folder root)
- /game         (your game client)
- /server       (your server program)
  - /infra
  - /deploy
- /profiles     (yml files containing deployment profiles)
- /resources    (yml files containing data for templates in game & server folders)
```


