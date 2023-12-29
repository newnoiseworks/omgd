# OMGD is Open Multiplayer Game Development

Multiplayer game dev is hard. OMGD aims to make some of it easier, namely, the setting up of servers on local machines for development and cloud instances for testing changes.

OMGD aims to be agnostic towards game engines and servers; the only requirement is that the game and servers can be built and ran using docker images and containers. For now, our examples focus on Godot as the game engine for client builds and dedicated servers, as well as using [Nakama](https://heroiclabs.com/nakama/) as a custom software server separate from the game build.

OMGD aims to be agnostic towards cloud providers via it's usage of [Terraform](https://www.terraform.io/), but for now only supports [Google Cloud Platform (GCP)](https://cloud.google.com/gcp).

OMGD at this time does NOT support production deployments, but aims to in the future.

## Full documentation

Full docs are available at [https://newnoiseworks.github.io/omgd](https://newnoiseworks.github.io/omgd)

### Install

For now, you need to grab a [binary from our release page](https://github.com/newnoiseworks/omgd/releases) for your OS and CPU setup, and put it on your computer's PATH so you can access it from the command line.

If you're on Windows, [here's a guide to adding executables to your PATH](https://stackoverflow.com/questions/4822400/register-an-exe-so-you-can-run-it-from-any-command-line-in-windows), and if you're on Mac / Linux, [here's a guide](https://medium.com/codex/adding-executable-program-commands-to-the-path-variable-5e45f1bdf6ce) for you.

### Local Development

#### Prerequisites

For local development only, you'll need `docker` and the `docker compose` plugin. We recommend using [Docker for Desktop](https://www.docker.com/products/docker-desktop/). You don't need a docker account.

#### Local setup

Start by downloading one of the [examples using Godot](https://github.com/orgs/newnoiseworks/repositories?q=example&type=all&language=&sort=) from our repositories.

In the example's folder, run the following commands.

```console
> omgd build-templates
> omgd game build
> omgd servers start
```

For details on each command and what it does, check the [full documentation](https://newnoiseworks.github.io/omgd).

Check your `game/dist` folder, there should be executables for the built game that you can load. If the `profiles/omgd.yml` file references a "web" servers target, you likely have a web based example built and running on `https://localhost:3030`. Also, try loading the project in your game enging from the `game` folder; buiding and running the game from the engine should work against the spawned server.

##### Check server logs

To see a running log from the servers, use the following command:

```console
> omgd servers logs
```

### Cloud Deployment for Testing

#### Prerequisites

For cloud deployment, you'll need to install [Terraform](https://www.terraform.io).

You'll also need a [Google Cloud Platform (GCP)](https://cloud.google.com/gcp) account with billing enabled.

<b>NOTE: OMGD is NOT responsible for your billing or finances on GCP or any cloud provider. You must monitor this. However, we default instances to `f1-micro`, and if you shut down your instances after testing, you shouldn't be paying more than a few pennies to dollars per month, USD.</b>

#### Cloud Setup

Please see the [getting started section](https://newnoiseworks.github.io/omgd/pages/getting-started/#deploying-to-the-cloud) for more.


