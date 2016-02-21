#!/usr/bin/env bash
GITROOT=`git rev-parse --show-toplevel`

export GOPATH="$GITROOT/gocode"
mkdir -p $GOPATH

mkdir -p $GOPATH/src/github.com/kurrik/opengl-benchmarks
rm -rf $GOPATH/src/github.com/kurrik/opengl-benchmarks/common
ln -s $GITROOT/src/common $GOPATH/src/github.com/kurrik/opengl-benchmarks/common

go get github.com/go-gl/gl/v3.3-core/gl
go get github.com/go-gl/glfw/v3.1/glfw
go get github.com/go-gl/mathgl/mgl32
go get github.com/golang/freetype
go get github.com/golang/freetype/truetype
go get golang.org/x/image/math/fixed

go run $GITROOT/src/uniform-texture-coords/*.go
