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

package sprites

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"unsafe"
)

type Sheet struct {
	keys    map[string]*Sprite
	texture *common.Texture
	Count   int
}

func NewSheet() *Sheet {
	return &Sheet{
		keys: map[string]*Sprite{},
	}
}

func (s *Sheet) SetTexture(texture *common.Texture) {
	s.Delete()
	s.texture = texture
}

func (s *Sheet) Bind() {
	if s.texture != nil {
		s.texture.Bind()
	}
}

func (s *Sheet) Unbind() {
	if s.texture != nil {
		s.texture.Unbind()
	}
}

func (s *Sheet) Delete() {
	if s.texture != nil {
		s.texture.Delete()
		s.texture = nil
	}
}

func (s *Sheet) AddSprite(key string, bounds, offset mgl32.Vec2) (out *Sprite) {
	var index int
	index = s.Count
	out = &Sprite{
		index:  index,
		bounds: bounds,
		offset: offset,
	}
	s.keys[key] = out
	s.Count++
	return
}

func (s *Sheet) Exists(key string) (exists bool) {
	_, exists = s.keys[key]
	return
}

func (s *Sheet) Sprite(key string) (out *Sprite, err error) {
	var exists bool
	if out, exists = s.keys[key]; !exists {
		err = fmt.Errorf("Invalid tile key %v", key)
		return
	}
	return
}

func (s *Sheet) Upload(ubo *common.UniformBuffer) (err error) {
	var (
		sprite *Sprite
		entry  uniformSprite
		data   = make([]uniformSprite, s.Count)
		size   = s.Count * int(unsafe.Sizeof(entry))
	)
	if s.texture == nil {
		err = fmt.Errorf("No texture associated with sheet")
		return
	}
	for _, sprite = range s.keys {
		data[sprite.index] = sprite.textureBounds(s.texture.Size)
	}
	ubo.Upload(data, size)
	return
}

func (s *Sheet) Texture() *common.Texture {
	return s.texture
}
