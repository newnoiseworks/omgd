!#/bin/bash

mkdir -v -p dist/windows
godot -v --export "Windows Desktop" dist/windows/game-${BUILD_ENV}.exe