!#/bin/bash

mkdir -v -p dist/web-${BUILD_ENV}
godot -v --export "HTML5" dist/web-${BUILD_ENV}/index.html