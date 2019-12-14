#!/bin/bash

mkdir -p build

set -x
for arch in 386 amd64 arm; do
  for os in darwin windows linux; do
    GOARCH=$arch GOOS=$os go build -o build/gorganizer-$os-$arch gorganizer.go
  done
done
set +x
