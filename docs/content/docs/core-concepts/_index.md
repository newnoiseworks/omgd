---
title: Core Concepts
---

# Core Concepts

OMGD is agnostic to game engines and servers so long as they can use docker containers to build the game, run the servers, and Terraform to setup server instances on the cloud.

To accomplish this, it requires a few things. Namely, a file structure to be adhered to containing `Dockerfile`s and `docker-compose.yml` files in certain locations, yml profiles to contain information to a certain format relating to those files, and templates to tie all the data together inside your game and server builds. 

## OMGD Workflow

{{<mermaid>}}
stateDiagram-v2
    Profiles: YML Profiles in the profiles/ dir
	note left of Profiles
		omgd infra commands supply profiles with IP address data of generated cloud VM instances
	end note
    Profiles --> Templates
    Templates: Profiles are used to provide data to the omgd build-templates command
	Templates --> BuiltTemplates
	BuiltTemplates: Templates produce config files for game build and servers matching profile data
    BuiltTemplates --> DockerDirectory
    DockerDirectory: Dockerfile and docker-compose.yml files are used by omgd to build and run game and game servers
    DockerDirectory -->  Terraform
    Terraform: For cloud deployments, Terraform and shell scripts are used to deploy to target systems
{{</mermaid>}}

## Concepts in Depth

### Debugging

OMGD can be debugged using env variables.

[Read about Debugging](/docs/core-concepts/debugging)

### File Structure

OMGD relies on a file structure which makes it agnostic to game engines and servers.

[Read about the File Structure](/docs/core-concepts/file-structure)

### Profiles

YML profiles allow for multiple environments and deployments of the game such to differ versions and demo changes independently.

[Read about the YML profiles](/docs/core-concepts/profiles)

### Templates

To feed information and keep them the same between game builds and servers, OMGD relies on a templating solution.

[Read about the Templates](/docs/core-concepts/templates)
