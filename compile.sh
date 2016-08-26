#!/bin/bash
# Cross-compile for container
env GOOS=linux GOARCH=amd64 go build
