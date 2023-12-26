---
title: OMGD
bookToc: false
# bookFlatSection: false
# bookHidden: false
# bookCollapseSection: false
# bookComments: false
# bookSearchExclude: false
---
# OMGD

## OMGD == Open Multiplayer GameDev

OMGD eases the entanglement of multiplayer game dev by combining open source tools so you can focus on what matters - your game.

All via a simple CLI tool and a few commands!
 
## How it works

This is how to start the local server and build game clients:

```tpl
 $ omgd build-templates
 $ omgd game build
 $ omgd servers start
```
<br />

This is how to setup a live game server - dedicated, another set of software, or both - onto a cloud service <small>(currently supporting Google Cloud Platform)</small>:

```tpl
 $ omgd infra instance-setup -p profiles/staging.yml
 $ omgd build-templates -p profiles/staging.yml
 $ omgd game build -p profiles/staging.yml
 $ omgd servers deploy -p profiles/staging.yml
```
<br />

OMGD uses Docker and Terraform for you to make everything work.

## Check out some examples

[Godot 3 Dedicated Server Example](https://github.com/newnoiseworks/omgd-godot3-dedicated-example)

## Try out the latest build

[Get the Latest Beta Release Today!](https://github.com/newnoiseworks/omgd/releases)

{{<hint warning>}}
In it's current stage OMGD is meant for development and testing stages of multiplayer game development only. Launching your game for production is on you, for now.
{{</hint>}}
