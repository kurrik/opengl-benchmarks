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

package common

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Buffer interface {
}

type GLBuffer struct {
	ID          uint32
	bufferBytes int
	target      uint32
}

func NewGLBuffer(target uint32) (b *GLBuffer) {
	b = &GLBuffer{
		target: target,
	}
	gl.GenBuffers(1, &b.ID)
	b.Bind()
	return
}

func (b *GLBuffer) Bind() {
	gl.BindBuffer(b.target, b.ID)
}

func (b *GLBuffer) Delete() {
	gl.DeleteBuffers(1, &b.ID)
}

func (b *GLBuffer) Upload(data interface{}, size int) {
	if size > b.bufferBytes {
		b.bufferBytes = size
		gl.BufferData(b.target, size, gl.Ptr(data), gl.STREAM_DRAW)
	} else {
		gl.BufferSubData(b.target, 0, size, gl.Ptr(data))
	}
}

type UniformBuffer struct {
	*GLBuffer
	binding uint32
}

func NewUniformBuffer() (b *UniformBuffer) {
	b = &UniformBuffer{
		GLBuffer: NewGLBuffer(gl.UNIFORM_BUFFER),
	}
	return
}

func (b *UniformBuffer) BlockBinding(programID uint32, name string, binding uint32) {
	b.binding = binding
	var (
		nameStr = gl.Str(fmt.Sprintf("%v\x00", name))
		index   = uint32(gl.GetUniformBlockIndex(programID, nameStr))
	)
	gl.UniformBlockBinding(programID, index, b.binding)
}

func (b *UniformBuffer) Upload(data interface{}, size int) {
	b.GLBuffer.Upload(data, size)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, b.binding, b.GLBuffer.ID, 0, size)
}

type ArrayBuffer struct {
	*GLBuffer
}

func NewArrayBuffer() (b *ArrayBuffer) {
	b = &ArrayBuffer{
		GLBuffer: NewGLBuffer(gl.ARRAY_BUFFER),
	}
	return
}

func (b *ArrayBuffer) VertexAttrib(programID uint32, name string, size int32, xtype uint32, stride int32, offset uintptr, divisor uint32) {
	var (
		nameStr   = gl.Str(fmt.Sprintf("%v\x00", name))
		location  = uint32(gl.GetAttribLocation(programID, nameStr))
		offsetPtr = gl.PtrOffset(int(offset))
	)
	gl.EnableVertexAttribArray(location)
	gl.VertexAttribPointer(location, size, xtype, false, stride, offsetPtr)
	gl.VertexAttribDivisor(location, divisor)
}
