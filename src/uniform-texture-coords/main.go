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
	"github.com/kurrik/opengl-benchmarks/common/binpacking"
	"github.com/kurrik/opengl-benchmarks/common/renderers"
	"image/color"
	"image/draw"
	"runtime"
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
		context       *common.Context
		sprites       *common.Sprites
		camera        *common.Camera
		framerate     *renderers.Framerate
		text          *renderers.Text
		font          *common.FontFace
		fg            = color.RGBA{255, 255, 255, 255}
		bg            = color.RGBA{0, 0, 0, 255}
		img           draw.Image
		packed        *binpacking.PackedImage
		packedTexture *common.Texture
		err           error
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
	if framerate, err = renderers.NewFramerateRenderer(); err != nil {
		panic(err)
	}
	if text, err = renderers.NewTextRenderer(); err != nil {
		panic(err)
	}
	if camera, err = context.Camera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{6, 4, 2}); err != nil {
		panic(err)
	}
	if font, err = common.NewFontFace("src/resources/Roboto-Light.ttf", 30, fg, bg); err != nil {
		panic(err)
	}
	if img, err = font.GetImage("foo bar baz Baj george"); err != nil {
		panic(err)
	}

	packed = binpacking.NewPackedImage(1024, 512)
	for _, s := range []string{
		"another string to add",
		"More string!",
		"Packing more string",
		"Add",
		"More",
		"String",
		"Framerate 59",
		"Framerate 60",
		"Framerate 20",
	} {
		if img, err = font.GetImage(s); err != nil {
			panic(err)
		}
		packed.Pack(s, img)
	}
	if err = common.WritePNG("test-packed.png", packed.Image()); err != nil {
		panic(err)
	}
	if err = common.WritePNG("test-font.png", img); err != nil {
		panic(err)
	}
	if packedTexture, err = common.GetTexture(packed.Image(), common.SmoothingLinear); err != nil {
		panic(err)
	}
	fmt.Printf("Sheet: %v\n", sprites.Sheet)
	for !context.ShouldClose() {
		context.Events.Poll()
		framerate.Bind()
		framerate.Render(camera)
		framerate.Unbind()
		packedTexture.Bind()
		text.Bind()
		text.Render(camera)
		text.Unbind()
		packedTexture.Unbind()
		context.SwapBuffers()
	}
}
