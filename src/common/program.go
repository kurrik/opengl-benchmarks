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
	"strings"
)

type Program struct {
	vao     uint32
	program uint32
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) Delete() {
	gl.DeleteVertexArrays(1, &p.vao)
	p.vao = 0
	gl.DeleteProgram(p.program)
	p.program = 0
}

func (p *Program) Bind() {
	gl.BindVertexArray(p.vao)
	gl.UseProgram(p.program)
}

func (p *Program) Unbind() {
	gl.BindVertexArray(0)
}

func (p *Program) Id() uint32 {
	return p.program
}

func (p *Program) Load(vertex, fragment string) (err error) {
	if err = p.createVAO(); err != nil {
		return
	}
	if err = p.buildProgram(vertex, fragment); err != nil {
		return
	}
	return
}

func (p *Program) createVAO() error {
	gl.GenVertexArrays(1, &p.vao)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR gl.GenVertexArray %X", e)
	}
	gl.BindVertexArray(p.vao)
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("ERROR array.Bind %X", e)
	}
	return nil
}

func (p *Program) compileShader(stype uint32, source string) (shader uint32, err error) {
	csource := gl.Str(source)
	shader = gl.CreateShader(stype)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)
	var status int32
	if gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status); status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))
		err = fmt.Errorf("ERROR shader compile:\n%s", log)
	}
	return
}

func (p *Program) linkProgram(vertex uint32, fragment uint32) (program uint32, err error) {
	program = gl.CreateProgram()
	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.BindFragDataLocation(program, 0, gl.Str("v_FragData\x00"))
	if e := gl.GetError(); e != 0 {
		err = fmt.Errorf("ERROR program.BindFragDataLocation %X", e)
		return
	}
	gl.LinkProgram(program)
	var status int32
	if gl.GetProgramiv(program, gl.LINK_STATUS, &status); status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		log := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))
		err = fmt.Errorf("ERROR program link:\n%s", log)
	}
	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)
	return
}

func (p *Program) buildProgram(vsrc string, fsrc string) (err error) {
	var (
		vertex   uint32
		fragment uint32
	)
	if vertex, err = p.compileShader(gl.VERTEX_SHADER, vsrc); err != nil {
		return
	}
	if fragment, err = p.compileShader(gl.FRAGMENT_SHADER, fsrc); err != nil {
		return
	}
	p.program, err = p.linkProgram(vertex, fragment)
	return
}
