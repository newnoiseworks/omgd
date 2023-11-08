# OMGD Bare Project Example

The purpose of the bare project is to set a structure used by all OMGD projects, universal to game engines and server backends.

It is not intended to work on it's own, and is just an example structure to build other templates on top of and around.

If OMGD does not support your chosen game engine and backend, you can use this to make a project OMGD can work with.

### Requirements

1. The folder structure and existing files must be present
2. The `game` directory needs a `DOCKERFILE` capable of building the game using a dockerized image of the game engine
3. The `game` directory also needs a `docker-compose.yml` file with services intended to build the game per each OS target. The `omgd.game.targets` block in the `omgd.yml` and other profiles will use the name of the services to allow you to select which OS' you support for your game.
4. The `servers` directory requires a backend e.g. Heroic Labs' Nakama to work properly, these files should be located in `servers/central`, and this folder will be uploaded recursively to the cloud instances per profile
5. The `servers` directory requiers a `docker-compose.yml` file to spin up containers per each service. This will be uploaded and used on your cloud instances as well, so make sure the folder paths are relative to only files located in the `servers` directory itself.
