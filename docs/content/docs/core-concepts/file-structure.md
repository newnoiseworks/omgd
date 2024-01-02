---
title: "File Structure"
---

# File Structure

OMGD relies on a file structure to operate. This assumes a few things:

1. You only have one game build, meaning your game is built in a single game engine and doesn't rely on multiple clients.
1. Dedicated game servers are a version of your client game build and can be built via the same build in the same game engine.
1. Your game can be built via a docker image defined by a `Dockerfile` and containers defined in a `docker-compose.yml` file.
1. You may have multiple custom game servers on top of a potential dedicated game server, for things like a web client, or a Nakama based server, or any other server technology that can run in a docker container.

The following is a breakdown of the required files and structure. Everything else is up to you. The [examples](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) listed show how this works a with Godot and other server types, dedicated and non.

```
- game/
  - Dockerfile
  - docker-compose.yml
- servers/
  - docker-compose.yml
- profiles/
  - omgd.yml
  - omgd.cloud.yml (for cloud deployments)
  - local.yml
  - your-profiles.yml
- resources/
  - your-custom-resources-for-templates.yml (not required)
```

## Game Directory

This is where your game engine will load it's project from. OMGD is agnostic to whatever that looks like so long as this is the root folder of the game project itself. OMGD requires two files to be in this directory to operate.

### Dockerfile

A `Dockerfile` is required to have an image which can build the game executables for you. This allows OMGD to operate like a CI/CD tool, making game builds automated. See the [examples](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) for what this looks like.

### docker-compose.yml

A `docker-compose.yml` file is required to use the above image and create containers which would generate game builds separated by OS target. This allows OMGD to build for various operating systems in parallel. Furthermore, the `omgd.game.targets` block of every profile needs to point to a `build-service` which would map to the services listed. See the [examples](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) for what this looks like.

## Servers Directory

This directory contains a folder per each server you wish to run. In development, local and deployed to a cloud, technically these servers will all run in a monolith style, meaning all within the same machine. In an example running a web client, a dedicated server, and a custom game server like Nakama, you should see three separate folders. See the [examples](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) for what this looks like.

### docker-compose.yml

This file is used to define, start, and stop servers within docker containers. It may point to a `Dockerfile`s within each server itself.

## Profiles Directory

This contains the YML profiles used to store base information on your game, cloud information to deploy your game, and then your own custom profiles to separate deployments as you see fit to different targets. This allows you to deploy a set of code to one target, change the code for development or testing purposes, and deploy the altered set to a new target, all with a yml file.

To read more about each individual file here, please see the [profiles section](/docs/core-concepts/profiles).

## Resources Directory

It can be useful to have game data stored in YML files that do not pertain to cloud deployments that are shared between your game builds and your servers. For an example, you may wish to have a store whose prices are linked between deployed game builds and servers, so they have the same data and can use this to prevent hacking. Here you can create YML files that will be referenceable in your templates.

To read more about resources and templates, see the [templates section](/docs/core-concepts/templates).
