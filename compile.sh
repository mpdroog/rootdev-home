#!/bin/bash
# Cross-compile for container
set +x
env GOOS=linux GOARCH=amd64 go build
docker build . -t "mpdroog/rootdev"
docker push mpdroog/rootdev
