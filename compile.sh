#!/bin/bash
# Cross-compile for container
set +x
env GOOS=linux GOARCH=amd64 go build
# Set ELF ABI to 003 to run under FreeBSD
echo -n $'\003' | dd bs=1 count=1 seek=7 conv=notrunc of=./rootdev

cd static-src
node_modules/gulp/bin/gulp.js
cd -

docker build . -t "mpdroog/rootdev"
docker push mpdroog/rootdev
