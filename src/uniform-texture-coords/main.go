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
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/renderers"
	"github.com/kurrik/opengl-benchmarks/common/spritesheet"
	"github.com/kurrik/opengl-benchmarks/common/text"
	"image/color"
	"runtime"
)

type Inst struct {
	Text string
	X    float32
	Y    float32
	R    float32
}

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	const (
		WinTitle  = "uniform-texture-coords"
		WinWidth  = 600
		WinHeight = 400
	)

	var (
		context   *common.Context
		sprites   *spritesheet.Sprites
		camera    *common.Camera
		framerate *renderers.Framerate
		font      *text.FontFace
		fg        = color.RGBA{255, 255, 255, 255}
		bg        = color.RGBA{0, 0, 0, 255}
		textMgr   *text.Manager
		err       error
		id        text.ID
		inst      *text.Instance
		rot       int = 0
	)
	if context, err = common.NewContext(); err != nil {
		panic(err)
	}
	if err = context.CreateWindow(WinWidth, WinHeight, WinTitle); err != nil {
		panic(err)
	}
	if sprites, err = spritesheet.NewSprites("src/resources/spritesheet.json", 32); err != nil {
		panic(err)
	}
	if framerate, err = renderers.NewFramerateRenderer(); err != nil {
		panic(err)
	}
	if textMgr, err = text.NewManager(100); err != nil {
		panic(err)
	}
	if camera, err = context.Camera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{6.4, 4.8, 2}); err != nil {
		panic(err)
	}
	if font, err = text.NewFontFace("src/resources/Roboto-Light.ttf", 24, fg, bg); err != nil {
		panic(err)
	}
	for _, s := range []Inst{
		Inst{Text: "This is text!", X: 0.05, Y: 0.05, R: 0},
		Inst{Text: "More text!", X: 1, Y: 1, R: 15},
	} {
		if id, err = textMgr.CreateText(); err != nil {
			panic(err)
		}
		if err = textMgr.SetText(id, s.Text, font); err != nil {
			panic(err)
		}
		if inst, err = textMgr.GetInstance(id); err != nil {
			return
		}
		inst.SetPosition(mgl32.Vec3{s.X, s.Y, 0})
		inst.SetRotation(s.R)
	}
	fmt.Printf("Sheet: %v\n", sprites.Sheet)
	for !context.ShouldClose() {
		context.Events.Poll()
		context.Clear()
		framerate.Bind()
		framerate.Render(camera)
		framerate.Unbind()
		textMgr.Bind()
		textMgr.Render(camera)
		textMgr.Unbind()
		context.SwapBuffers()

		textMgr.SetText(id, fmt.Sprintf("Rotation %v", rot%10), font)
		inst.SetRotation(float32(rot))
		rot += 1
	}
	if err = common.WritePNG("test-packed.png", textMgr.PackedImage.Image()); err != nil {
		panic(err)
	}
	textMgr.Delete()
}
