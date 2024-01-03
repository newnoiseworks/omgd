---
title: "omgd servers"
---

## `> omgd servers`

Manages docker containers that run all necessary servers to run your game. Essentially runs docker commands against provided profiles and game builds.

{{<hint warning>}}
With the exception of `deploy`, the subcommands below don't work on remote servers defined in profiles when using OMGD on Windows. If you wish to run them, for now you'll need to ssh into the remote servers yourself. All commands should work on local machines though.
{{</hint>}}

## Subcommands

### `> omgd servers deploy`

Usable after creating a deployment target via `omgd infra instance-setup -p profiles/your-profile.yml`

Requires a non local profile to be referenced with a valid `omgd.servers.host` value to deploy to. Will copy the entire `servers/` folder as it exists to the remote server - be sure to run `omgd game build -p profiles/your-profile.yml` before hand to make sure it uploads builds matching this profile. It will then start docker containers running each server.

```
omgd servers deploy -p profiles/your-profile.yml
```

Works on Windows.

### `> omgd servers start`

Starts docker containers representing your servers. Useful in local development. Usable remotely but the above `deploy` command runs the servers independently on their own. Uses the local profile by default. If a non local profile is supplied, it will try to run the command on a remote server.

```
omgd servers start
```

If using a non local profile, this command doesn't work on Windows against remote servers for now.

### `> omgd servers logs`

Shows a running log of statements coming from all docker containers. Requires docker containers to be running. Uses the local profile by default. If a non local profile is supplied, it will try to run the command on a remote server.

```
omgd servers logs
```

If using a non local profile, this command doesn't work on Windows against remote servers for now.

### `> omgd servers stop`

Stops docker containers representing your servers. Useful in local development to preserve resources or to restart the server after rebuilding it. Uses the local profile by default. If a non local profile is supplied, it will try to run the comand on a remote server.

```
omgd servers stop
```

If using a non local profile, this command doesn't work on Windows against remote servers for now.

See below flag for dropping data volumes w/ `omgd servers stop`

## Flags

### Drop volumes (`omgd servers stop -v`)

Only usable with `omgd servers stop`. If any of your servers are writing data to a shared volume created by your docker containers, this will drop that data volume. Useful for database instances like `PostgreSQL` running via Docker for servers like `Nakama` et al.

```
omgd servers stop -v
```
