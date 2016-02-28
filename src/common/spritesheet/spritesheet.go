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

package spritesheet

import (
	"github.com/go-gl/mathgl/mgl32"
)

type SpritesheetFrame struct {
	Size          mgl32.Vec2 // In units
	TextureOffset mgl32.Vec2
	TextureSize   mgl32.Vec2
}

type Spritesheet struct {
	frames      map[string]*SpritesheetFrame
	texturePath string
}

func NewSpritesheet(texturePath string) *Spritesheet {
	return &Spritesheet{
		frames:      map[string]*SpritesheetFrame{},
		texturePath: texturePath,
	}
}

func (s Spritesheet) GetTexturePath() string {
	return s.texturePath
}

func (s Spritesheet) GetFrame(name string) *SpritesheetFrame {
	if _, present := s.frames[name]; !present {
		return nil
	}
	return s.frames[name]
}

func (s *Spritesheet) AddFrame(name string, frame *SpritesheetFrame) {
	s.frames[name] = frame
}
