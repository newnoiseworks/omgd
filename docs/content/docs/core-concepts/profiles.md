---
title: "Profiles"
---

# Profiles

For most internet driven applications, games or otherwise, having multiple servers representing different states of your application or game is necessary.

For example, most apps or games have a running "production" server or set of servers that is used by the public consumers of the app or game. Internally, the developers have their own set of servers meant to test their work. Many organizations have numerous "staging" or "development" servers for this, and a process to test, merge, and deploy changes to the "production" server when necessary to address bugs and new features.

OMGD's profiles allow for this with multiplayer game dev by using profiles that are meant to represent different deployment targets.

{{<hint warning>}}
OMGD is not built for production use at this time. Please make your own considerations when setting up your game for production and only use OMGD for development and testing, for now.
{{</hint>}}

## YML Syntax

OMGD profiles are `.yml` files stored in the `/profiles` folder at the root directory of the OMGD project. [YML, or YAML](https://yaml.org/) is a syntax meant to store data, similar to JSON and XML. It's fairly straightforward. Check out the `/profiles` directory in any of the [examples](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) to see some `.yml` files.

## Profile inheritance pattern

OMGD Profiles follow an inheritance pattern. There is one top level profile, or the parent profile, called `omgd.yml`. It stores the general information on your game and is used by all profiles as a base of information to adjust depending on what each profile needs.

Afterwards, if a `omgd.cloud.yml` file exists, it will then add that after the `omgd.yml` file. The difference is that `omgd.cloud.yml` stores all data on your cloud services provider. OMGD currently only supports GCP.

Finally, the profile you wish to use against `omgd` CLI commands is used. By default a `local.yml` file is always available for local development and is used by `omgd` commands if no other profile is supplied. The `local.yml` file is meant to be only used for local development on the user's machine.

If you make a new `yml` file, for example `staging.yml`, that is intended to map to a remote compute instance. Note the filename is used by `omgd` to name the compute instance.

{{<hint danger>}}
If you create a yml file, then create an instance using `omgd infra instance-setup`, renaming the file can lose the link between the instance and `omgd`. After creating an instance, it's important to not rename the files.
{{</hint>}}

{{<mermaid>}}
stateDiagram-v2
	omgdyml: omgd.yml top level profile
    omgdyml --> omgdcloudyml
    omgdcloudyml: omgd.cloud.yml file optionally used for basic cloud configuration
	omgdcloudyml --> deploytargetyml
	deploytargetyml: your-profile.yml file used to define deployment target information and override omgd.yml settings as needed
{{</mermaid>}}

## OMGD YML Profile Reference

These values are used to configure `omgd.yml` and any custom profiles e.g. `local.yml` as well.

A hashtag `#` denotes a code comment.

```
omgd:
  name: your-project-name
  game: 
    targets:
      - build-service: build-windows # references game/docker-compose.yml service
      - build-service: build-x11-server
        copy: "game/dist/x11-server" # OMGD can copy game builds, useful for web and dedicated servers
        to: "servers/dedicated-build/src"
  servers:
    ports: # list all TCP and UDP ports used by your servers for cloud setup
      tcp: '"22", "3030", "7348-7351", "9999"'
      udp: '"7348-7351", "9999"'
    services:
      - build-service: web # references servers/docker-compose.yml service
      - build-service: dedicated
      - build-service: nakama
```

## OMGD Cloud YML Profile Reference

The [examples](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=) all have a `profiles/omgd.cloud.example.yml` file to build off of. Rename it to `omgd.cloud.yml` and edit the information to reference your cloud config setup. Currently only GCP is supported.

```
omgd:
  gcp:
    project: your-project-id-on-gcp
    zone: your-gcp-zone # e.g. us-east-4c etc
    creds-file: override/default/gcp/creds-file.json # optional
```

{{<hint info>}}
`omgd.gcp.creds-file` is optional and can be used if you have a json file reprsenting a GCP service account that is in a location other than the default json created by the `gcloud auth application-default login` command.
{{</hint>}}

## Overrides

All the above keys can be overriden in the profiles you create as well as `local.yml` for your own needs. For instance, you may want to deploy a test without building all of the clients built in the `omgd.yml` file, or maybe more. Or, maybe you want to change the cloud zone that you are deploying to, if say one developer would rather deploy servers to the west coast than the east coast in America.
