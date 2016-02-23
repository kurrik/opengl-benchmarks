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
	vbo           uint32
	vboBytes      int
	ubo           uint32
	uboBytes      int
	uboBinding    uint32
	stride        int32
	locColor      int32
	locModelView  int32
	locProjection int32
	data          *textData
}

func NewTextRenderer() (r *Text, err error) {
	r = &Text{
		shader: common.NewProgram(),
	}
	if err = r.shader.Load(TEXT_VERTEX, TEXT_FRAGMENT); err != nil {
		return
	}
	r.shader.Bind()
	gl.GenBuffers(1, &r.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.GenBuffers(1, &r.ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, r.ubo)
	var (
		point            textDataPoint
		locWorldPosition = uint32(gl.GetAttribLocation(r.shader.Id(), gl.Str("v_WorldPosition\x00")))
		offWorldPosition = gl.PtrOffset(int(unsafe.Offsetof(point.worldPos)))
		locTile          = uint32(gl.GetAttribLocation(r.shader.Id(), gl.Str("f_Tile\x00")))
		offTile          = gl.PtrOffset(int(unsafe.Offsetof(point.tile)))
		uboIndex         = uint32(gl.GetUniformBlockIndex(r.shader.Id(), gl.Str("TextureData\x00")))
	)
	r.uboBinding = 1
	gl.UniformBlockBinding(r.shader.Id(), uboIndex, r.uboBinding)
	r.stride = int32(unsafe.Sizeof(point))
	r.locColor = gl.GetUniformLocation(r.shader.Id(), gl.Str("v_Color\x00"))
	r.locModelView = gl.GetUniformLocation(r.shader.Id(), gl.Str("m_ModelView\x00"))
	r.locProjection = gl.GetUniformLocation(r.shader.Id(), gl.Str("m_Projection\x00"))
	gl.EnableVertexAttribArray(locWorldPosition)
	gl.EnableVertexAttribArray(locTile)
	gl.VertexAttribPointer(locWorldPosition, 2, gl.FLOAT, false, r.stride, offWorldPosition)
	gl.VertexAttribDivisor(locWorldPosition, 1)
	gl.VertexAttribPointer(locTile, 1, gl.FLOAT, false, r.stride, offTile)
	gl.VertexAttribDivisor(locTile, 1)
	return
}

func (r *Text) Bind() {
	r.shader.Bind()
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, r.ubo)
}

func (r *Text) Unbind() {
	r.shader.Unbind()
}

func (r *Text) Delete() {
	r.shader.Delete()
	// TODO: Delete UBO and VBO
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
		modelView     = mgl32.Ident4()
		dataBytes int = len(r.data.Points) * int(r.stride)
	)
	gl.Uniform4f(r.locColor, 0, 255.0/255.0, 0, 255.0/255.0)
	gl.UniformMatrix4fv(r.locModelView, 1, false, &modelView[0])
	gl.UniformMatrix4fv(r.locProjection, 1, false, &camera.Projection[0])

	if dataBytes > r.vboBytes {
		r.vboBytes = dataBytes
		gl.BufferData(gl.ARRAY_BUFFER, dataBytes, gl.Ptr(r.data.Points), gl.STREAM_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, dataBytes, gl.Ptr(r.data.Points))
	}

	textureData := []float32{1, 1, 0, 0}
	uboBytes := int(unsafe.Sizeof(textureData))
	if uboBytes > r.uboBytes {
		r.uboBytes = uboBytes
		gl.BufferData(gl.UNIFORM_BUFFER, uboBytes, gl.Ptr(textureData), gl.STREAM_DRAW)
	} else {
		gl.BufferSubData(gl.UNIFORM_BUFFER, 0, uboBytes, gl.Ptr(textureData))
	}
	gl.BindBufferRange(gl.UNIFORM_BUFFER, r.uboBinding, r.ubo, 0, uboBytes)

	//gl.DrawArrays(gl.POINTS, 0, int32(len(r.data.Points)))
	ptsPerInstance := 6
	instanceCount := 2
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, int32(ptsPerInstance), int32(instanceCount))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}
