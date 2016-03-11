// Copyright 2016 Arne Roomann-Kurrik
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

package render

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"unsafe"
)

type Point struct {
	Position mgl32.Vec3
	Texture  mgl32.Vec2
	Frame    float32
}

type Geometry struct {
	Points []Point
	Dirty  bool
	vbo    *common.ArrayBuffer
	stride uintptr
}

func NewGeometry(capacity int) (out *Geometry) {
	var (
		point  Point
		stride uintptr = unsafe.Sizeof(point)
	)
	out = &Geometry{
		Points: make([]Point, 0, capacity),
		Dirty:  true,
		stride: stride,
		vbo:    common.NewArrayBuffer(stride),
	}
	return
}

func (g *Geometry) Register(shaderID uint32, posName, texName, tileName string) {
	var point Point
	g.vbo.Vec3(shaderID, posName, unsafe.Offsetof(point.Position), 0)
	g.vbo.Vec2(shaderID, texName, unsafe.Offsetof(point.Texture), 0)
	g.vbo.Float(shaderID, tileName, unsafe.Offsetof(point.Frame), 0)
}

func (g *Geometry) Bind() {
	g.vbo.Bind()
}

func (g *Geometry) Delete() {
	if g.vbo != nil {
		g.vbo.Delete()
		g.vbo = nil
	}
}

func (g *Geometry) Upload() {
	if g.Dirty {
		g.vbo.Upload(g.Points, len(g.Points)*int(g.stride))
		g.Dirty = false
	}
}
