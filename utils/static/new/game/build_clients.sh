#!/bin/bash

BUILD_ENV=$BUILD_ENV docker-compose run -T build-mac
docker rmi -f newnoiseworks/game-build-mac-$BUILD_ENV
BUILD_ENV=$BUILD_ENV docker-compose run -T build-windows
docker rmi -f newnoiseworks/game-build-windows-$BUILD_ENV
BUILD_ENV=$BUILD_ENV docker-compose run -T build-web
docker rmi -f newnoiseworks/game-build-web-$BUILD_ENV
BUILD_ENV=$BUILD_ENV docker-compose run -T build-x11
docker rmi -f newnoiseworks/game-build-x11-$BUILD_ENV
