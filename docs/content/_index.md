---
title: Multiplayer Game Dev Made Open and Easy - OMGD
addDot: false
addBang: true
---

OMGD eases the entanglement of multiplayer game dev by combining open source tools so you can focus on what matters - your game.

All via a simple CLI tool and a few commands!
<br />
<br />
<br />
 
## How it works

This is how to start the local server and build game clients:

```

 $ omgd build-templates
 $ omgd build-clients
 $ omgd server start
 
```

This is how to setup and teardown a live game server - dedicated, another set of software, or both - onto a cloud service <small>(currently supporting Google Cloud Platform)</small>:

```

 $ omgd infra deploy -p staging.yml
 $ omgd infra game-deploy -p staging.yml
 $ omgd infra destroy -p staging.yml
 
```

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
