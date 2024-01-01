---
title: "omgd build-templates"
---

## `> omgd build-templates`

This will search through your directories for all files ending in `.tmpl`. It will use the profile provided (`profiles/local.yml` by default) to supply data to that file, generating a file which strips the final `.tmpl` from it's name, e.g. `example.gd.tmpl` will create a file `example.gd`.

These files are used to inject IP address and hostname data of newly created servers into the game and server builds so multiple deployments can easily exist and be driven by your YML profile definitions.

See the [templates core concept page](/docs/core-concepts/templates) for more.

## General Form

`omgd build-templates [OPTIONS]`

## Options

### Extension (`--ext [STRING]`)

By default, the `build-templates` command will search for files with the extension `tmpl`. You can adjust this by supplying your own extension. Do not include the `.` before the extension name, just the letters itself.

E.g. `omgd build-tempaltes --ext customtmpl`

### Remove (`--remove`)

If you pass this flag, the template files will be deleted after generating. Useful in edge cases when using OMGD with other build and deployment systems.

e.g. `omgd build-templates --remove`
