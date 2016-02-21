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
in vec2 v_TexturePositionOut;
uniform vec4 v_Color;
out vec4 v_FragData;
void main() {
  v_FragData = v_Color;
}` + "\x00"

const TEXT_VERTEX = `#version 150
in vec2 v_WorldPosition;
in vec2 v_TexturePosition;
uniform mat4 m_ModelView;
uniform mat4 m_Projection;
out vec2 v_TexturePositionOut;
void main() {
  gl_Position = m_Projection * m_ModelView * vec4(v_WorldPosition, 0.0, 1.0);
  v_TexturePositionOut = v_TexturePosition;
}` + "\x00"

type textDataPoint struct {
	worldPos   mgl32.Vec2
	texturePos mgl32.Vec2
}

type textData struct {
	Points []textDataPoint
}

type Text struct {
	shader        *common.Program
	vbo           uint32
	vboBytes      int
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
	var (
		point              textDataPoint
		locWorldPosition   = uint32(gl.GetAttribLocation(r.shader.Id(), gl.Str("v_WorldPosition\x00")))
		offWorldPosition   = gl.PtrOffset(int(unsafe.Offsetof(point.worldPos)))
		locTexturePosition = uint32(gl.GetAttribLocation(r.shader.Id(), gl.Str("v_TexturePosition\x00")))
		offTexturePosition = gl.PtrOffset(int(unsafe.Offsetof(point.texturePos)))
	)
	r.stride = int32(unsafe.Sizeof(point))
	r.locColor = gl.GetUniformLocation(r.shader.Id(), gl.Str("v_Color\x00"))
	r.locModelView = gl.GetUniformLocation(r.shader.Id(), gl.Str("m_ModelView\x00"))
	r.locProjection = gl.GetUniformLocation(r.shader.Id(), gl.Str("m_Projection\x00"))
	gl.EnableVertexAttribArray(locWorldPosition)
	gl.EnableVertexAttribArray(locTexturePosition)
	gl.VertexAttribPointer(locWorldPosition, 2, gl.FLOAT, false, r.stride, offWorldPosition)
	gl.VertexAttribPointer(locTexturePosition, 2, gl.FLOAT, false, r.stride, offTexturePosition)
	return
}

func (r *Text) Bind() {
	r.shader.Bind()
	gl.BindBuffer(gl.ARRAY_BUFFER, r.vbo)
}

func (r *Text) Unbind() {
	r.shader.Unbind()
}

func (r *Text) Delete() {
	r.shader.Delete()
}

func (r *Text) Render(camera *common.Camera) (err error) {
	// Temporary:
	r.data = &textData{
		Points: []textDataPoint{
			textDataPoint{
				worldPos: mgl32.Vec2{0, 0},
				texturePos:  mgl32.Vec2{0, 0},
			},
			textDataPoint{
				worldPos: mgl32.Vec2{0.25, 0.25},
				texturePos:  mgl32.Vec2{0.25, 0.25},
			},
			textDataPoint{
				worldPos: mgl32.Vec2{0.5, 0.5},
				texturePos:  mgl32.Vec2{0.5, 0.5},
			},
			textDataPoint{
				worldPos: mgl32.Vec2{0.75, 0.75},
				texturePos:  mgl32.Vec2{0.75, 0.75},
			},
			textDataPoint{
				worldPos: mgl32.Vec2{1, 1},
				texturePos:  mgl32.Vec2{1, 1},
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

	gl.DrawArrays(gl.POINTS, 0, int32(len(r.data.Points)))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR: OpenGL error %X", e)
	}
	return
}
