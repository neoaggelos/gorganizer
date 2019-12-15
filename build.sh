#!/bin/bash

version="0.1"

mkdir -p build
set -x
for arch in 386 amd64 arm; do
  for os in darwin windows linux; do
    [[ "x$os" = "xwindows" ]] && ext=".exe" || ext=""
    GOARCH=$arch GOOS=$os go build -o "build/gorganizer-${version}-${os}-${arch}${ext}" gorganizer.go
  done
done
set +x
