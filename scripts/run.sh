#!/usr/bin/env bash
GITROOT=`git rev-parse --show-toplevel`

export GOPATH="$GITROOT/gocode"
mkdir -p $GOPATH

go get github.com/go-gl/gl/v3.3-core/gl
go get github.com/go-gl/glfw/v3.1/glfw
go get github.com/go-gl/mathgl/mgl32

go run $GITROOT/src/uniform-texture-coords/*.go
