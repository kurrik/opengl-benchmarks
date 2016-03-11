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
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/sheet"
	"unsafe"
)

const FRAGMENT = `#version 150

precision mediump float;

in vec2 v_TexturePos;
in vec2 v_TextureMin;
in vec2 v_TextureDim;
uniform sampler2D u_Texture;
out vec4 v_FragData;

void main() {
  vec2 v_TexturePosition = v_TextureMin + mod(v_TexturePos, v_TextureDim);
  v_FragData = texture(u_Texture, v_TexturePosition);
}`

const VERTEX = `#version 150

#define MAX_TILES 1024

struct Tile {
  vec4 texture;
};

layout (std140) uniform TextureData {
  Tile Tiles[MAX_TILES];
};

in vec3 v_Position;
in vec2 v_Texture;
in float f_Tile;
in mat4 m_Model;
uniform mat4 m_View;
uniform mat4 m_Projection;
out vec2 v_TexturePos;
out vec2 v_TextureMin;
out vec2 v_TextureDim;

void main() {
  Tile t_Tile = Tiles[int(f_Tile)];
  v_TextureMin = t_Tile.texture.zw;
  v_TextureDim = t_Tile.texture.xy;
  v_TexturePos = v_Texture * v_TextureDim;
  gl_Position = m_Projection * m_View * m_Model * vec4(v_Position, 1.0);
}`

type renderInstance struct {
	model mgl32.Mat4
	frame float32
}

type Renderer struct {
	shader     *common.Program
	vbo        *common.ArrayBuffer
	ubo        *common.UniformBuffer
	uView      *common.Uniform
	uProj      *common.Uniform
	bufferSize int
	buffer     []renderInstance
	stride     uintptr
}

func NewRenderer(bufferSize int) (r *Renderer, err error) {
	var (
		point  renderInstance
		stride = unsafe.Sizeof(point)
	)
	r = &Renderer{
		shader:     common.NewProgram(),
		bufferSize: bufferSize,
		buffer:     make([]renderInstance, bufferSize),
		stride:     stride,
	}
	if err = r.shader.Load(VERTEX, FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()
	r.vbo = common.NewArrayBuffer(stride)
	r.vbo.Float(r.shader.ID(), "f_Frame", unsafe.Offsetof(point.frame), 1)
	r.vbo.Mat4(r.shader.ID(), "m_Model", unsafe.Offsetof(point.model), 1)
	r.ubo = common.NewUniformBuffer(r.shader.ID())
	r.ubo.BlockBinding("TextureData", 1)
	r.uView = r.shader.Uniform("m_View")
	r.uProj = r.shader.Uniform("m_Projection")
	return
}

func (r *Renderer) Bind() {
	r.shader.Bind()
	r.ubo.Bind()
	r.vbo.Bind()
}

func (r *Renderer) Register(geometry *Geometry) {
	geometry.Bind()
	geometry.Register(r.shader.ID(), "v_Position", "v_Texture", "f_Tile")
}

func (r *Renderer) Unbind() {
	r.shader.Unbind()
}

func (r *Renderer) Delete() {
	if r.shader != nil {
		r.shader.Delete()
		r.shader = nil
	}
	if r.ubo != nil {
		r.ubo.Delete()
		r.ubo = nil
	}
	if r.vbo != nil {
		r.vbo.Delete()
		r.vbo = nil
	}
}

func (r *Renderer) draw(geometry *Geometry, count int) (err error) {
	if count <= 0 {
		return
	}
	r.vbo.Upload(r.buffer, count*int(r.stride))
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(len(geometry.Points)), int32(count))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}

func (r *Renderer) Render(
	camera *common.Camera,
	regions sheet.UniformBufferRegions,
	geometry *Geometry,
	instances *InstanceList,
) (err error) {
	var (
		instance *Instance
		i        *renderInstance
		index    int
	)
	r.uView.Mat4(camera.View)
	r.uProj.Mat4(camera.Projection)
	geometry.Bind()
	geometry.Upload()
	regions.Upload(r.ubo)
	index = 0
	instance = instances.Head()
	if instance != nil {
		if index >= r.bufferSize {
			if err = r.draw(geometry, index); err != nil {
				return
			}
			index = 0
		}
		i = &r.buffer[index]
		i.frame = float32(instance.Frame)
		i.model = instance.GetModel()
		index++
		instance = instance.Next()
	}
	err = r.draw(geometry, index)
	return
}
