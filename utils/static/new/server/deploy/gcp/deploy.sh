#!/bin/bash

truncate -s 0 /var/log/docker.log

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $PWD:$PWD -w=$PWD docker/compose:1.27.0 down

docker run --rm -v /var/run/docker.sock:/var/run/docker.sock -v $PWD:$PWD -w=$PWD docker/compose:1.27.0 up -d

