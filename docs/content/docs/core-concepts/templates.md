---
title: "Templates"
---

# Templates

OMGD uses a templating solution driven by [Golang Templates](https://pkg.go.dev/text/template) to supply server information (IP addresses, host names, and more) per each deployment target.

`omgd build-templates` looks through the `game` and `servers` directories for `.tmpl` files. In most of the [examples using Godot game engine](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) you should see a `game/Autoloads/GameConfig.gd.tmpl` file. Running `omgd build-templates` should produce a `game/Autoloads/GameConfig.gd` file against the local profile so your game build knows the IP address of your servers, in the case of `local.yml`, this will generally be `127.0.0.1`.

## Profile data

The profile you pass into the `omgd build-templates` call (`local.yml` if none is supplied) can be referenced by prepending the word "profile" to the value in the template.

e.g. `{{.profile.omgd.name}}`

## Resources data

The `/resources` dir stores optional YML files that can be used to feed data that is game and non profile specific, for example, a version stamp, into your game and servers. You can store whatever you'd like, and this will be referenceable in any `.tmpl` files you create.

For instance, if you have a `/resources/game.yml` file that looks like the following:

```
game:
  version: 1.2
```

This can feed data into a `GameVersion.gd.tmpl` file like so:

```
var version = {{.game.version}}
```

If you create another file in your `servers/` directory somewhere in a running server, this way you can map the versions so that after deployment you can ensure the game and servers match version via your own code.

## Template Functions

OMGD supports a few custom functions should they be useful to you in your templates.

### MD5

OMGD can convert a string into an md5 hash.

e.g. `{{ md5 .your.resource.value }}`

### Upper Snake

OMGD can convert a string to be in "UPPER_SNAKE" format.

e.g. `{{ upperSnake .your.resource.value }}`

### Camel case

OMGD can convert a string to "camelCase" format.

e.g. `{{ camel .your.resource.value }}`
