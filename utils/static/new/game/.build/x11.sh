!#/bin/bash

mkdir -v -p dist/x11
godot -v --export "Linux/X11" dist/x11/game-${BUILD_ENV}.x86_64