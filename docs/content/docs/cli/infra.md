---
title: "omgd infra"
---

## `> omgd infra`

Manages all cloud infrastructure creation for your project via internal Terraform files and deployment scripts.

If you wish to use your own Terraform and deployment approach, you must create your own files and run them manually in your own processes. You can use [`omgd build-templates`](/docs/cli/build-templates) to inject profile data into terraform and deployment related scripts as needed.

{{<hint info>}}
Commands will copy Terraform and deployment scripts to a hidden `.omgd` folder in your project. At the end of the command's run, unless an error occured, this folder will be deleted.
{{</hint>}}

## Subcommands

### `> omgd infra project-setup`

This sets up the project on your cloud provider. This will create a data bucket to store terraform config files per each deployment, as well as a VPS network with firewall rules aligned to your port configurations in the profile.

This only needs to be run once against a project. Note that it does NOT set up game servers, it is just to create the resources necessary to do so.

Note that this does not require a profile, and will use the data found in `omgd.yml` as well as `omgd.cloud.yml` for all resource creation.

### `> omgd infra project-destroy`

This will destroy the data bucket and VPS network created in `omgd infra project-setup`.

Note that if this cannot run, you may manually delete the project on your cloud infrastructure, or it's individually created GCS bucket and VPS network which should be named against the `omgd.name` value in your top level `omgd.yml` profile.

### `> omgd infra instance-setup -p [PROFILE]`

e.g. `omgd infra instance-setup -p profiles/staging.yml`

This will setup a VM (virtual machine) on your cloud infrastructure that you can deploy built game servers to. Note that it does NOT deploy the game server, just sets up the cloud infrastructure necessary to deploy them later via [`omgd servers deploy`](/docs/cli/servers).

Note that it requires a non local profile. It will write in the newly created VM instance IP into the `omgd.servers.host` value of the provided profile.

### `> omgd infra instance-destroy -p [PROFILE]`

e.g. `omgd infra instance-destroy -p profiles/staging.yml`

This will destroy any VMs setup via a previously run `instance-setup` command. It should reset any `omgd.servers.host` values to `???` after destroying the VM.
