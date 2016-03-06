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

package batch

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"unsafe"
)

type rInstance struct {
	tile  float32
	point mgl32.Vec3
}

const BATCH_FRAGMENT = `#version 150

precision mediump float;

in vec2 v_TexturePos;
in vec2 v_TextureMin;
in vec2 v_TextureDim;
uniform sampler2D u_Texture;
out vec4 v_FragData;

void main() {
  vec2 v_TexturePosition = v_TextureMin + mod(v_TexturePos - v_TextureMin, v_TextureDim);
  v_FragData = texture(u_Texture, v_TexturePosition);
}`

const BATCH_VERTEX = `#version 150

#define MAX_TILES 1024

struct Tile {
  vec4 texture;
  vec4 size;
};

layout (std140) uniform TextureData {
  Tile Tiles[MAX_TILES];
};

in float f_Tile;
in vec3 v_Point;
uniform mat4 m_Model;
uniform mat4 m_View;
uniform mat4 m_Projection;
out vec2 v_TexturePos;
out vec2 v_TextureMin;
out vec2 v_TextureDim;

void main() {
  Tile t_Tile = Tiles[int(f_Tile)];
  v_TexturePos = v_Point.xy / 4;
  v_TextureMin = t_Tile.texture.zw;
  v_TextureDim = t_Tile.texture.xy;
  gl_Position = m_Projection * m_View * m_Model * vec4(v_Point, 1.0);
}`

type Renderer struct {
	shader *common.Program
	stride uintptr
	ubo    *common.UniformBuffer
	vbo    *common.ArrayBuffer
	uView  *common.Uniform
	uProj  *common.Uniform
	uModel *common.Uniform
}

func NewRenderer() (r *Renderer, err error) {
	r = &Renderer{
		shader: common.NewProgram(),
	}
	if err = r.shader.Load(BATCH_VERTEX, BATCH_FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()
	var point rInstance
	r.stride = unsafe.Sizeof(point)
	r.vbo = common.NewArrayBuffer(r.shader.ID(), r.stride)
	r.vbo.Vec3("v_Point", unsafe.Offsetof(point.point), 0)
	r.vbo.Float("f_Tile", unsafe.Offsetof(point.tile), 0)
	r.ubo = common.NewUniformBuffer(r.shader.ID())
	r.ubo.BlockBinding("TextureData", 1)
	r.uView = r.shader.Uniform("m_View")
	r.uProj = r.shader.Uniform("m_Projection")
	r.uModel = r.shader.Uniform("m_Model")
	return
}

func (r *Renderer) Bind() {
	r.shader.Bind()
	r.vbo.Bind()
	r.ubo.Bind()
}

func (r *Renderer) Unbind() {
	r.shader.Unbind()
}

func (r *Renderer) Delete() {
	r.shader.Delete()
	r.vbo.Delete()
	r.ubo.Delete()
}

func (r *Renderer) Render(camera *common.Camera, sheet *tile.Sheet) (err error) {
	r.uModel.Mat4(mgl32.Ident4())
	r.uView.Mat4(camera.View)
	r.uProj.Mat4(camera.Projection)
	r.ubo.Upload(sheet.Tiles, sheet.TileBytes())
	data := []rInstance{
		rInstance{tile: 4, point: mgl32.Vec3{0, 0, 0}},
		rInstance{tile: 4, point: mgl32.Vec3{1, 0, 0}},
		rInstance{tile: 4, point: mgl32.Vec3{1, 1, 0}},
	}
	r.vbo.Upload(data, len(data)*int(r.stride))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(data)))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}
