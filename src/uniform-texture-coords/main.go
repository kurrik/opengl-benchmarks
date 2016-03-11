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
	"flag"
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/glog"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/batch"
	"github.com/kurrik/opengl-benchmarks/common/render"
	"github.com/kurrik/opengl-benchmarks/common/spritesheet"
	"github.com/kurrik/opengl-benchmarks/common/text"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"github.com/kurrik/opengl-benchmarks/common/util"
	"image/color"
	"runtime"
)

const BATCH = `
AAA
BBB
`

type Inst struct {
	Key string
	X   float32
	Y   float32
	R   float32
}

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

func main() {
	flag.Parse()

	const (
		WinTitle              = "uniform-texture-coords"
		WinWidth              = 640
		WinHeight             = 480
		PixelsPerUnit float32 = 100
	)

	var (
		context       *common.Context
		spriteMgr     *spritesheet.Manager
		camera        *common.Camera
		framerate     *util.Framerate
		font          *text.FontFace
		fg            = color.RGBA{255, 255, 255, 255}
		bg            = color.RGBA{0, 0, 0, 255}
		textMgr       *text.Manager
		err           error
		inst          *tile.Instance
		rot           int = 0
		batchRenderer *batch.Renderer
		textMapping   *batch.TextMapping
		textLoader    *batch.TextLoader
		batchData     *batch.Batch
		renderer      *render.Renderer
	)
	if context, err = common.NewContext(); err != nil {
		panic(err)
	}
	if err = context.CreateWindow(WinWidth, WinHeight, WinTitle); err != nil {
		panic(err)
	}
	if spriteMgr, err = spritesheet.NewManager(spritesheet.Config{
		JsonPath:      "src/resources/spritesheet.json",
		PixelsPerUnit: PixelsPerUnit,
		MaxInstances:  100,
	}); err != nil {
		panic(err)
	}
	if framerate, err = util.NewFramerateRenderer(); err != nil {
		panic(err)
	}
	if textMgr, err = text.NewManager(text.Config{
		MaxInstances:  100,
		TextureWidth:  512,
		TextureHeight: 512,
		PixelsPerUnit: PixelsPerUnit,
	}); err != nil {
		panic(err)
	}
	if camera, err = context.Camera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{6.4, 4.8, 2}); err != nil {
		panic(err)
	}
	if font, err = text.NewFontFace("src/resources/Roboto-Light.ttf", 24, fg, bg); err != nil {
		panic(err)
	}
	for _, s := range []Inst{
		Inst{Key: "This is text!", X: 0, Y: -1.0, R: 0},
		Inst{Key: "More text!", X: 1.0, Y: 1.0, R: 15},
	} {
		if inst, err = textMgr.CreateInstance(); err != nil {
			panic(err)
		}
		if err = textMgr.SetText(inst, s.Key, font); err != nil {
			panic(err)
		}
		inst.SetPosition(mgl32.Vec3{s.X, s.Y, 0})
		inst.SetRotation(s.R)
	}

	for _, s := range []Inst{
		Inst{Key: "numbered_squares_02", X: 0, Y: 0, R: 0},
		Inst{Key: "numbered_squares_02", X: -2.0, Y: -2.0, R: -15},
	} {
		if inst, err = spriteMgr.CreateInstance(); err != nil {
			panic(err)
		}
		if err = spriteMgr.SetFrame(inst, s.Key); err != nil {
			panic(err)
		}
		inst.SetPosition(mgl32.Vec3{s.X, s.Y, 0})
		inst.SetRotation(s.R)
	}

	if batchRenderer, err = batch.NewRenderer(); err != nil {
		panic(err)
	}
	if textMapping, err = batch.NewTextMapping(spriteMgr.Regions(), "numbered_squares_03"); err != nil {
		panic(err)
	}
	textMapping.Set('A', "numbered_squares_01")
	textMapping.Set('B', "numbered_squares_tall_16")
	textLoader = batch.NewTextLoader(textMapping)
	if batchData, err = textLoader.Load(batchRenderer, 1, BATCH); err != nil {
		panic(err)
	}

	if renderer, err = render.NewRenderer(100); err != nil {
		panic(err)
	}

	//fmt.Printf("Sheet: %v\n", sprites.Tiles)
	for !context.ShouldClose() {
		context.Events.Poll()
		context.Clear()

		batchRenderer.Bind()
		spriteMgr.Regions().Texture().Bind()
		batchRenderer.Render(camera, spriteMgr.Regions(), batchData)
		batchRenderer.Unbind()

		spriteMgr.Bind()
		spriteMgr.Render(camera)
		spriteMgr.Unbind()

		framerate.Bind()
		framerate.Render(camera)
		framerate.Unbind()

		textMgr.Bind()
		textMgr.Render(camera)
		textMgr.Unbind()

		renderer.Bind()
		//renderer.Render(camera, spriteMgr.Regions(), ..., ...)
		renderer.Unbind()

		context.SwapBuffers()

		if err = textMgr.SetText(textMgr.Instances.Head(), fmt.Sprintf("Rotation %v", rot%100), font); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			break
		}
		inst.SetRotation(float32(rot))
		rot += 1
	}
	if err = common.WritePNG("test-packed.png", textMgr.Regions().Image()); err != nil {
		panic(err)
	}
	textMgr.Delete()
	glog.Flush()
}
