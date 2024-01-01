---
title: CLI API
---

# OMGD CLI API

OMGD is a command line interface based application. Below are the basic commands and their descriptions. Click the command for deeper information and options.

## Command List

| <div style="min-width:200px">Command</div> | Description |
| -- | -- |
| [build-templates](/docs/cli/build-templates) | Builds template files using OMGD profile data. |
| [game](/docs/cli/game) | Builds game projects via docker. |
| [infra](/docs/cli/infra) | Sets up cloud infrastructure and development instances via terraform and deployment scripts. |
| [servers](/docs/cli/servers) | Manages docker containers, most useful for local development, limited support for commands against cloud instances if supplied a profile attached to a cloud instance. |

## Global Flags

| <div style="min-width:150px">Option Flag</div> | Default | Description |
| -- | -- | -- |
| `-p`, `--profile` | `profiles/local.yml` | The OMGD YML profile to run command against. See [profiles section](/docs/core-concepts/profiles) for more details on profiles. |
| `--output-dir` | `.` (current working directory) | The directory that OMGD should output files into. Defaults to current directory. |
