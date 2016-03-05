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

package text

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"unsafe"
)

const TEXT_FRAGMENT = `#version 150
precision mediump float;
in vec2 v_TexturePosition;
uniform sampler2D u_Texture;
out vec4 v_FragData;
void main() {
  v_FragData = texture(u_Texture, v_TexturePosition);
}` + "\x00"

const TEXT_VERTEX = `#version 150

#define MAX_TILES 1024

struct Tile {
  vec4 texture;
  vec4 size;
};

layout (std140) uniform TextureData {
  Tile Tiles[MAX_TILES];
};

const vec2 WorldPoints[] = vec2[6](
  vec2(-0.5, -0.5),
  vec2( 0.5,  0.5),
  vec2(-0.5,  0.5),
  vec2(-0.5, -0.5),
  vec2( 0.5, -0.5),
  vec2( 0.5,  0.5)
);

const vec2 TexturePoints[] = vec2[6](
  vec2(0.0, 0.0),
  vec2(1.0, 1.0),
  vec2(0.0, 1.0),
  vec2(0.0, 0.0),
  vec2(1.0, 0.0),
  vec2(1.0, 1.0)
);

in float f_Tile;
in mat4 m_Model;
uniform mat4 m_View;
uniform mat4 m_Projection;
out vec2 v_TexturePosition;
void main() {
  Tile t_Tile = Tiles[int(f_Tile)];
  mat4 m_ModelView = m_View * m_Model;
  vec4 v_Point = vec4(WorldPoints[gl_VertexID] * t_Tile.size.xy, 0.0, 1.0);
  gl_Position = m_Projection * m_ModelView * v_Point;
  v_TexturePosition = TexturePoints[gl_VertexID] * t_Tile.texture.xy + t_Tile.texture.zw;
}` + "\x00"

type Renderer struct {
	shader *common.Program
	stride uintptr
	data   *rendererData
	ubo    *common.UniformBuffer
	vbo    *common.ArrayBuffer
	uView  *common.Uniform
	uProj  *common.Uniform
}

func NewRenderer() (r *Renderer, err error) {
	r = &Renderer{
		shader: common.NewProgram(),
	}
	if err = r.shader.Load(TEXT_VERTEX, TEXT_FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()
	var point rendererInstance
	r.stride = unsafe.Sizeof(point)
	r.vbo = common.NewArrayBuffer(r.shader.ID(), r.stride)
	r.vbo.Float("f_Tile", unsafe.Offsetof(point.tile), 1)
	r.vbo.Mat4("m_Model", unsafe.Offsetof(point.model), 1)
	r.ubo = common.NewUniformBuffer(r.shader.ID())
	r.ubo.BlockBinding("TextureData", 1)
	r.uView = r.shader.Uniform("m_View")
	r.uProj = r.shader.Uniform("m_Projection")
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

func (r *Renderer) Render(camera *common.Camera, data *rendererData, textureData *tile.Sheet) (err error) {
	var (
		vboBytes = data.Count * int(r.stride)
		uboBytes = textureData.TileBytes()
	)
	r.uView.Mat4(camera.View)
	r.uProj.Mat4(camera.Projection)
	r.vbo.Upload(data.Instances, vboBytes)
	r.ubo.Upload(textureData.Tiles, uboBytes)
	ptsPerInstance := 6
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(ptsPerInstance), int32(data.Count))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}
