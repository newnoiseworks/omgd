---
title: 'Getting Started'
menu:
  before:
    weight: 1
---

# Getting Started with OMGD

## Install

Grab the latest release [from GitHub](https://github.com/newnoiseworks/omgd/releases) and make sure it's executable from your operating system's PATH.

If you're on Windows, [here's a guide to adding executables to your PATH](https://stackoverflow.com/questions/4822400/register-an-exe-so-you-can-run-it-from-any-command-line-in-windows), and if you're on Mac / Linux, [here's a guide](https://medium.com/codex/adding-executable-program-commands-to-the-path-variable-5e45f1bdf6ce) for you.

## Local Development

To get started, grab one of the examples from [GitHub](https://github.com/orgs/newnoiseworks/repositories?q=omgd+example&type=all&language=&sort=).

### Requirements

For local development you'll need the following installed:

1. Docker with Docker Compose (we recommend Docker for Desktop)

### Building the Game and Servers

Then, in the directory of the project on your computer, run

```console
omgd build-templates
omgd game build
omgd servers start
```

In the `game/dist` directory, there should be built executables mapped to your local server. Launch a couple to try the demo out!

Also, generally loading the `game` directory within the game engine of the example should let you build the game and it will work against the server created.

{{<hint info>}}
`omgd game build` may take awhile - it's building the game into executables using docker containers and the image defined in `game/Dockerfile`. For more verbose output, check the [debugging section](/docs/core-concepts/debugging).
{{</hint>}}


### Checking Server Logs

To verify running servers, check the logs.

```console
omgd servers logs
```

### Shutting Down the Servers

To stop the servers, useful if you need to rebuild the game's web client or dedicated server using `omgd game build` and start the servers again with the new code, just run the following.

```console
omgd servers stop 
```

## Deploying to the Cloud

### Requirements

For deploying to the cloud, you'll need the following installed:

1. Terraform
1. GCloud CLI SDK

{{<hint warning>}}
You will need billing enabled on your GCP account. The costs for servers can be pennies if you stick to f1-micro compute instances (the default) and take them down when they're not being tested on. OMGD is not responsible for managing your billing.
{{</hint>}}

### GCloud SDK Setup

For deploying to the Google Cloud Platform (GCP), ensure you're logged into your GCloud account via the CLI so Terraform can work.

```console
gcloud auth application-default login
```

Via the GCP website do the following:

1. Make sure you have billing enabled on your GCP account - creating a new one comes with free credits but requires a credit card.
1. Create a new project and take note of the project ID. You'll need it in a bit.
1. For that project, enable the Compute Engine API (search for "compute engine" in the search bar after you've created the project).

### OMGD Cloud Settings Setup

Copy the `profiles/omgd.cloud.example.yml` file in your profiles directory, rename it to `profiles/omgd.cloud.yml`

In the `omgd.gcp.project` field, put in the project ID of your GCP project you recently setup.

Then, run the project setup command:

```console
omgd infra project-setup
```

After it's done, reload your `omgd.cloud.yml` file and you should see a value for the property `omgd.gcp.bucket` - your cloud settings will be saved to a newly created GCS bucket.

{{<hint info>}}
This command may take some time.
{{</hint>}}

### Create a profile to deploy your game

Each `profile/*.yml` file not named `local.yml`, `omgd.yml`, or `omgd.cloud.yml` creates a new instance to deploy your game servers to while you're in development.

Do the following to deploy:

1. Create a new file named `profiles/example.yml`. You can leave it blank for now, or you can use it to override values setup in `profiles/omgd.yml`, the parent profile, if you wish.
2. Run the following command: 

```console
omgd infra instance-setup -p profiles/example.yml
```

After awhile it should finish and your `profiles/example.yml` should have an IP with your newly created instance!

{{<hint danger>}}
The filename of your profile is used when creating servers. If you create a server using the above command, and then change the filename, they will lose the connection and you will need to manage the server manually.
{{</hint>}}

### Build game for cloud deployment

Your game needs to know the IP address (and other info) of your new profile.

Build the game against that profile with the following:

```console
omgd build-templates -p profiles/example.yml
omgd game build -p profiles/example.yml
```

Note that these commands are the same as the above for local development, just with a `-p` flag to set the profile. Without it, the `profiles/local.yml` profile is used.

{{<hint info>}}
When you're ready to get back to local development, you may need to run `omgd build-templates` or `omgd game build` so your game and servers are readied to run on your local machine.
{{</hint>}}

### Deploy game servers to the cloud

As easy as:

```console
omgd servers deploy -p profiles/example.yml
```

### Cleaning up the server

First destroy the instance with the following:

```console
omgd infra instance-destroy -p profiles/example.yml
```

Then, if you wish, destroy the project. This will disable new instances from being created.

```console
omgd infra project-destroy
```
