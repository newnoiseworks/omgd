---
title: 'Getting Started'
menu:
  before:
    weight: 1
image: images/writer.jpeg
---

# Getting Started with OMGD

## Requirements

For local development you'll need the following installed:

1. Docker with Docker Compose (we reccomend Docker for Desktop)

For deploying to the cloud, you'll need the following installed:

1. Terraform
1. GCloud CLI SDK

{{<hint warning>}}
You will need billing enabled on your GCP account. The costs for servers can be pennies if you stick to f1-micro compute instances (the default) and take them down when they're not being tested on.
{{</hint>}}

## Install

Grab the latest release [from Github](https://github.com/newnoiseworks/omgd/releases) and make sure it's readable from your operating system's PATH

## Local Development

To get started, grab one of the examples from [Github](https://github.com/newnoiseworks/).

### Building the Game and Starting the Servers

Then, in the directory of the project on your computer, run

```console
omgd game build
omgd server start
```

Follow the README.md in the example but generally loading the `game` directory within the game engine of the example should let you build the game and it will work against the server created. 

Also, in the `game/dist` directory, there should be built executables.

To verify, check the server logs.

### Checking Server Logs

```console
omgd server logs
```

## Deploying to the Cloud

### GCloud SDK Setup

For deploying to the Google Cloud Platform (GCP), ensure you're logged into your GCloud account via the cli so Terraform can work.

```console
gcloud auth application-default login
```

Via the GCP website do the following:

0. Make sure you have billing enabled on your GCP account - creating a new one comes with free credits but requires a credit card.
1. Create a new project and take note of the project ID. You'll need it in a bit.
2. For that project, enable the compute engine API.

### OMGD Cloud Settings Setup

Copy the `profiles/omgd.cloud.example.yml` file in your profiles directory, rename it to `profiles/omgd.cloud.yml`

In the `omgd.gcp.project` field, put in the project ID of your GCP project you recently setup.

Then, run the project setup command:

```console
omgd infra project-setup
```

After it's done, reload your `omgd.cloud.yml` file and you should see a value for the property `omgd.gcp.bucket` - your cloud settings will be saved to a newly created GCS bucket.

### Create a profile to deploy your game to

Each `profile/*.yml` file not named `local.yml`, `omgd.yml`, or `omgd.cloud.yml` creates a new instance to deploy your game servers to while you're in development.

Do the following to deploy:

1. Create a new file named `profiles/example.yml`. You can leave it blank for now, or you can use it to override values setup in `profiles/omgd.yml`, the parent profile, if you wish.
2. Run the following command: 

```console
omgd infra instance-setup -p profiles/example.yml
```

After awhile it should finish and your `profiles/example.yml` should have an IP with your newly created instance!

### Build game for cloud deployment

Your game needs to know the IP address (and other info) of your new profile.

Build the game against that profile with the following:

```console
omgd game build -p profiles/example.yml
```

### Deploy game servers to the cloud

As easy as:

```console
omgd servers deploy -p profiles/example.yml
```

