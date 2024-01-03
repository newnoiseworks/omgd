# OMGD is Open Multiplayer Game Development

Multiplayer game dev is hard. OMGD aims to make some of it easier, namely, the setting up of servers on local machines for development and cloud instances for testing changes.

OMGD aims to be agnostic towards game engines and servers; the only requirement is that the game and servers can be built and ran using docker images and containers.

For now, our examples focus on [Godot](https://godotengine.org) as the game engine for client builds and dedicated servers, as well as using [Nakama](https://heroiclabs.com/nakama/) as a custom software server separate from the game build.

OMGD aims to be agnostic towards cloud providers via it's usage of [Terraform](https://www.terraform.io/), but for now only supports [Google Cloud Platform (GCP)](https://cloud.google.com/gcp).

OMGD at this time does NOT support production deployments, but aims to in the future. OMGD is meant for development and testing purposes in the game dev cycle only.

## Full documentation

Full docs are available at [https://newnoiseworks.github.io/omgd](https://newnoiseworks.github.io/omgd)

### Install

For now, you need to grab a [binary from our release page](https://github.com/newnoiseworks/omgd/releases) for your OS and CPU setup, and put it on your computer's PATH so you can access it from the command line.

If you're on Windows, [here's a guide to adding executables to your PATH](https://stackoverflow.com/questions/4822400/register-an-exe-so-you-can-run-it-from-any-command-line-in-windows), and if you're on Mac / Linux, [here's a guide](https://medium.com/codex/adding-executable-program-commands-to-the-path-variable-5e45f1bdf6ce) for you.

### Getting Started

Please see the [getting started section](https://newnoiseworks.github.io/omgd/pages/getting-started/) for more.


