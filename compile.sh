#!/bin/bash
# Cross-compile for container
set +x
env GOOS=linux GOARCH=amd64 go build

cd static-src
node_modules/gulp/bin/gulp.js
cd -

docker build . -t "mpdroog/rootdev"
docker push mpdroog/rootdev
