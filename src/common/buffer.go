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
	"unsafe"
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
	programID uint32
	binding   uint32
}

func NewUniformBuffer(programID uint32) (b *UniformBuffer) {
	b = &UniformBuffer{
		programID: programID,
		GLBuffer:  NewGLBuffer(gl.UNIFORM_BUFFER),
	}
	return
}

func (b *UniformBuffer) BlockBinding(name string, binding uint32) {
	b.binding = binding
	var (
		nameStr = gl.Str(fmt.Sprintf("%v\x00", name))
		index   = uint32(gl.GetUniformBlockIndex(b.programID, nameStr))
	)
	gl.UniformBlockBinding(b.programID, index, b.binding)
}

func (b *UniformBuffer) Upload(data interface{}, size int) {
	b.GLBuffer.Upload(data, size)
	gl.BindBufferRange(gl.UNIFORM_BUFFER, b.binding, b.GLBuffer.ID, 0, size)
}

type ArrayBuffer struct {
	*GLBuffer
	programID uint32
	stride    uintptr
}

func NewArrayBuffer(programID uint32, stride uintptr) (b *ArrayBuffer) {
	b = &ArrayBuffer{
		GLBuffer:  NewGLBuffer(gl.ARRAY_BUFFER),
		programID: programID,
		stride:    stride,
	}
	return
}

func (b *ArrayBuffer) vertexAttrib(location uint32, size int32, xtype uint32, offset uintptr, divisor uint32) {
	fmt.Printf("LOCATION %v\n", location)
	var offsetPtr = gl.PtrOffset(int(offset))
	gl.EnableVertexAttribArray(location)
	gl.VertexAttribPointer(location, size, xtype, false, int32(b.stride), offsetPtr)
	gl.VertexAttribDivisor(location, divisor)
}

func (b *ArrayBuffer) VertexAttrib(name string, size int32, xtype uint32, offset uintptr, divisor uint32) {
	var (
		nameStr  = gl.Str(fmt.Sprintf("%v\x00", name))
		location = uint32(gl.GetAttribLocation(b.programID, nameStr))
	)
	b.vertexAttrib(location, size, xtype, offset, divisor)
}

func (b *ArrayBuffer) Float(name string, offset uintptr, divisor uint32) {
	b.VertexAttrib(name, 1, gl.FLOAT, offset, divisor)
}

func (b *ArrayBuffer) Vec2(name string, offset uintptr, divisor uint32) {
	b.VertexAttrib(name, 2, gl.FLOAT, offset, divisor)
}

func (b *ArrayBuffer) Mat4(name string, offset uintptr, divisor uint32) {
	var (
		float    float32
		sizeVec4 = unsafe.Sizeof(float) * 4
		nameStr  = gl.Str(fmt.Sprintf("%v\x00", name))
		location = uint32(gl.GetAttribLocation(b.programID, nameStr))
	)
	fmt.Printf("LOCATION %v\n", location)
	b.vertexAttrib(location, 4, gl.FLOAT, offset, divisor)
	b.vertexAttrib(location+1, 4, gl.FLOAT, offset+sizeVec4, divisor)
	b.vertexAttrib(location+2, 4, gl.FLOAT, offset+2*sizeVec4, divisor)
	b.vertexAttrib(location+3, 4, gl.FLOAT, offset+3*sizeVec4, divisor)
}
