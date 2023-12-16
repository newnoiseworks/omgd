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
<br />
<br />
<br />
 
## How it works

This is how to start the local server and build game clients:

```
 $ omgd game build
 $ omgd servers start
```
<br />

This is how to setup a live game server - dedicated, another set of software, or both - onto a cloud service <small>(currently supporting Google Cloud Platform)</small>:

```
 $ omgd infra instance-setup -p staging.yml
 $ omgd game build -p staging.yml
 $ omgd servers deploy -p staging.yml
```
<br />

OMGD uses Docker and Terraform for you to make everything work.
<br />
<br />
<br />

## Check out some examples

[Godot 3 Dedicated Server Example](https://github.com/newnoiseworks/omgd-godot3-dedicated-example)
<br />
<br />
<br />

## Try out the latest build

[Get the Latest Beta Release Today!](https://github.com/newnoiseworks/omgd/releases)
<br />
<br />
<br />
