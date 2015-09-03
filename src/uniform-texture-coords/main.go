// Copyright 2015 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"../common"
	"fmt"
	"runtime"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	const (
		WinTitle  = "uniform-texture-coords"
		WinWidth  = 640
		WinHeight = 480
	)

	var (
		context   *common.Context
		sprites   *common.Sprites
		camera    *common.Camera
		framerate *common.RendererFramerate
		err       error
	)
	if context, err = common.NewContext(); err != nil {
		panic(err)
	}
	if err = context.CreateWindow(WinWidth, WinHeight, WinTitle); err != nil {
		panic(err)
	}
	if sprites, err = common.NewSprites("src/resources/spritesheet.json", 32); err != nil {
		panic(err)
	}
	if framerate, err = common.NewRendererFramerate(); err != nil {
		panic(err)
	}
	if camera, err = context.Camera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{6, 4, 2}); err != nil {
		panic(err)
	}
	fmt.Printf("Sheet: %v\n", sprites.Sheet)
	for !context.ShouldClose() {
		context.Events.Poll()
		framerate.Bind()
		framerate.Render(camera)
		framerate.Unbind()
		context.SwapBuffers()
	}
}
