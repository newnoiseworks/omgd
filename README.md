# gg
### GameDev Glue

#### Automation by configuration for your multi tiered game development projects

## Vision

`gg` is a task runner meant to "glue" together your game's various projects -- infrastructure setup, game client, game servers, more, maybe less -- with your existing build and deploy processes into a single configuration file meant to reflect an environment -- your local computer, a demo server, production - whatever -- and store it in a repository for everyone on the team to share.

`gg` also adds resource backed templating, so you can share data between your projects and easily spawn several "versions" of your game all using different configurations for your testing needs.

See [`omgd`](https://github.com/newnoiseworks/omgd) for a project which uses the `gg` runner base in setting up projects for you with a preconfigured build, deploy, and infra setup process using Godot, Nakama, and Terraform.

## Installation

Currently you must download this repo and then build the binary, placing it in your path. For example, given `~/.local/bin` is in your `$PATH` variable, then the following would be good:

`$ go build -o ~/.local/bin/gg`

Eventually we'll post releases with version tags on this repo, and then after that via popular package managers like snap, brew etc.


## Usage
#### Warning: In development!

There are three things to gg: *profiles*, *resources*, and *templates*

Projects for now must be setup as follows at the root. See profiles/example.yml for an example profile.

```
/               (your project folder root)
- /profiles     (yml files containing deployment profiles)
- /resources    (yml files containing data for templates in game & server folders)
```

### Profiles:

Profiles are yml files held in the `/profiles` directory at the root of your project. `gg` interprets two fields, both arrays, one named `main` which is required and the other optionally named `tasks`.

`main` goes through each of it's items, called `project steps`, and executes it's internal `steps` one by one. The `dir` attribute can be used at both levels to control which directory the command runs on.

`tasks` are optional sub tasks you can setup.

#### An example profile yml file:

e.g. `profiles/staging.yml`
```yml
main:
  - name: infra
    dir: server/infra
    steps: 
      - cmd: gg build-templates .
      - cmd: ./infra_deploy_if_needed.sh
        dir: gcp # infers ./server/infra/gcp based on top level dir
  - name: setup-ip
    steps:
      - cmd: gg update-profile nakama.host $(./server/infra/get_ip.sh)
  - name: game-client
    dir: game
    steps: 
      - cmd: gg build-templates .
      - cmd: ./build_clients.sh
      - cmd: ./deploy_clients.sh
  - name: nakama-server
    dir: server
    steps: 
      - cmd: gg build-templates .
      - cmd: ./build_server.sh
      - cmd: ./deploy/gcp/deploy.sh


tasks:
  - name: destroy-infra
    steps:
      - cmd: ./destroy_infra.sh
    dir: server/infra

# Feel free to populate the below w/ whatever. 
# In your templates they will be namespaced under "profile"

example_server:
  host_ip: ??? # gets defined above in main tasks
```

#### gg Update profiles

You can make on-the-fly updates to your gg profile. This is useful to change data that may be updated, such as a host entry after an infrastructure update. Simply call

`$ gg update-profile value.key value --profile=profiles/example`

Note that you can omit `--profile=` if you're calling this from a profile setup, as in the example above.


### gg Resources

`gg` encourages you to store your _game_ specific content in the `resources` folder -- this is useful for templating and is explained below.


### gg Templating

One of `gg`'s main utilities is to add templating across all of your projects. This can be used for obvious things like tying together server configurations -- our example above updates this profile's "example_server.host_ip" entry so the subsequent game and server build tasks can use that value in their templates.

It can also be used for "game" configuration that involves the client and the server needing to know the same things -- like current prices on items -- and maybe in cases where the client needs a partial set of data and the server can control the rest -- all of that is stored in yml files in the `resources` folder.

So the `gg build-templates [dir] --profile=profiles/example` command does the following:

1. Searches your `resources` directory and combines the `yml` files together in one big object, and combines the current chosen profile (`profiles/local.yml` by default) as well so the resulting object has all your `resources` key value pairs as well as your profile under the key `profile` -- note that keys will be overriden should you have repetitions at the root level, so we reccomend each `resources` file leading with one top level "holder" key named after the file.
2. Looks through the given directory for files ending in `.tmpl`, and when found, write a new file to the same directory shedding the `.tmpl` in the filename and using the object in #1 as a data source for the template. The template uses [go's template solution](https://golang.org/pkg/text/template/) and syntax.

