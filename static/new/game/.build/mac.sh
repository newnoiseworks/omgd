!#/bin/bash

mkdir -v -p dist/mac
godot -v --export "Mac OSX" dist/mac/game-${BUILD_ENV}.zip