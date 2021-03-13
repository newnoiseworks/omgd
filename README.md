# TPL-Fred

#### A build and deployment tool for _The Promised Land_'s game client, server, and potentially (one day) it's website

## Initial project setup instructions

For now we'll be storing this info here though it's likely to move and this README will be exclusively focused on the build tool.

That said:

0. Pull this repository down into an empty parent directory for the project. E.g.:

```
TPLProjectFolder
  | tpl-fred
```

1. In the `build/profiles` folder of this library, make a copy of `example.yml` and call it `local.yml` -- edit it only if you altered your local nakama server config for some reason, but otherwise leave it alone.

2. Within this repo, run this command:
   `$ go run main.go clone local --output=../`

   This will clone the game and server repositories into the parent directory w/ appropriate names for this build tool.

3. In the same directory, run these commands:
   `$ go run main.go build-config game local --output=../`
   `$ go run main.go build-config server local --output=../`

   They setup the necessary build artifacts (item lists, mission info, more, see `builder/config/templates` within this repo for details) so the game and server can run.

4. Now, assuming you have docker running with the `docker-compose` CLI available:
   `$ cd ../server`
   `$ docker-compose up -d`

   This will stand up the nakama server locally as well as a local instance of the database it uses, CockroachDB. Try running `docker-compose ps` in the same repo to ensure they're both up stably.

5. Assuming that went well, open the freshly cloned `game` repository folder w/n Godot. Try running the game with the play button, and you should be able to log in.
