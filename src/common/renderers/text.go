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

package renderers

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
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

layout (std140) uniform TextureData {
  vec4 Tiles[MAX_TILES];
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
in vec2 v_WorldPosition;
uniform mat4 m_ModelView;
uniform mat4 m_Projection;
out vec2 v_TexturePosition;
void main() {
  vec4 Tile = Tiles[int(f_Tile)];
  v_TexturePosition = TexturePoints[gl_VertexID] * Tile.xy + Tile.zw;
  vec2 WorldPosition = WorldPoints[gl_VertexID] + v_WorldPosition;
  gl_Position = m_Projection * m_ModelView * vec4(WorldPosition, 0.0, 1.0);
}` + "\x00"

type textDataPoint struct {
	worldPos mgl32.Vec2
	tile     float32
}

type textData struct {
	Points []textDataPoint
}

type Text struct {
	shader        *common.Program
	stride        uintptr
	locColor      int32
	locModelView  int32
	locProjection int32
	data          *textData
	ubo           *common.UniformBuffer
	vbo           *common.ArrayBuffer
}

func NewTextRenderer() (r *Text, err error) {
	r = &Text{
		shader: common.NewProgram(),
	}
	if err = r.shader.Load(TEXT_VERTEX, TEXT_FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()
	var point textDataPoint
	r.stride = unsafe.Sizeof(point)
	r.vbo = common.NewArrayBuffer()
	r.vbo.VertexAttrib(r.shader.ID(), "v_WorldPosition", 2, gl.FLOAT, r.stride, unsafe.Offsetof(point.worldPos), 1)
	r.vbo.VertexAttrib(r.shader.ID(), "f_Tile", 1, gl.FLOAT, r.stride, unsafe.Offsetof(point.tile), 1)

	r.ubo = common.NewUniformBuffer()
	r.ubo.BlockBinding(r.shader.ID(), "TextureData", 1)

	r.locColor = gl.GetUniformLocation(r.shader.ID(), gl.Str("v_Color\x00"))
	r.locModelView = gl.GetUniformLocation(r.shader.ID(), gl.Str("m_ModelView\x00"))
	r.locProjection = gl.GetUniformLocation(r.shader.ID(), gl.Str("m_Projection\x00"))
	return
}

func (r *Text) Bind() {
	r.shader.Bind()
	r.vbo.Bind()
	r.ubo.Bind()

}

func (r *Text) Unbind() {
	r.shader.Unbind()
}

func (r *Text) Delete() {
	r.shader.Delete()
	r.vbo.Delete()
	r.ubo.Delete()
}

func (r *Text) Render(camera *common.Camera) (err error) {
	// Temporary:
	r.data = &textData{
		Points: []textDataPoint{
			textDataPoint{
				worldPos: mgl32.Vec2{0, 0},
				tile:     0,
			},
			textDataPoint{
				worldPos: mgl32.Vec2{1, 1},
				tile:     0,
			},
		},
	}
	var (
		modelView   = mgl32.Ident4()
		vboBytes    = len(r.data.Points) * int(r.stride)
		textureData = []float32{1, 1, 0, 0}
		uboBytes    = int(unsafe.Sizeof(textureData))
	)
	gl.Uniform4f(r.locColor, 0, 255.0/255.0, 0, 255.0/255.0)
	gl.UniformMatrix4fv(r.locModelView, 1, false, &modelView[0])
	gl.UniformMatrix4fv(r.locProjection, 1, false, &camera.Projection[0])
	r.vbo.Upload(r.data.Points, vboBytes)
	r.ubo.Upload(textureData, uboBytes)
	ptsPerInstance := 6
	instanceCount := 2
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(ptsPerInstance), int32(instanceCount))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}
