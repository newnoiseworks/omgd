---
title: "omgd game"
---

## `> omgd game`

This command is used for building and managing games in the `game/` folder. It will use the profile provided and the parent `omgd.yml` profile to infer what targets are run.

The targets that are defined in the profiles link to service names in the required `game/docker-compose.yml` file. In the example projects, these default to `build-[OPERATING SYSTEM NAME]`, e.g. `build-windows` for windows executables.

The profile listing for each target can also supply a directory to copy the built game to. This is useful for moving web client builds as well as dedicated server builds into appropriate folders in the `servers/` directories such to run the game as internet accessibly servers. See the [profiles core concept page](/docs/core-concepts/profiles) for more information.

## Subcommands

### `> omgd game build`

Builds the game against the supplied profile, `profiles/local.yml` by default. To build a game against a profile named staging.yml:

```
omgd game build -p profiles/staging.yml
```

To build against the default local profile, just run

```
omgd game build
```

## Flags

### Targets (`--targets [STRING]`)

Supply targets as named by the `game/docker-compose.yml` services to run game builds against. Can be useful to limit build targets for testing.

```
omgd build game --targets "build-x11 build-windows build-web"
```

This will ignore the `omgd.game.targets` section of the compiled profile and attempt to run `docker` commands against the supplied targets which should map to build services as specified in `game/docker-compose.yml`
