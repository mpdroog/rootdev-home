#!/bin/bash
# Cross-compile for container
set -e
set -x
env GOOS=linux GOARCH=amd64 go build
# Set ELF ABI to 003 to run under FreeBSD
echo -n $'\003' | dd bs=1 count=1 seek=7 conv=notrunc of=./rootdev

cd static-src
yarn
node_modules/gulp/bin/gulp.js
cd -

COMMIT=$(git rev-parse HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
docker build . -t "mpdroog/rootdev" --build-arg GIT_COMMIT="$BRANCH-$COMMIT"
docker push mpdroog/rootdev
