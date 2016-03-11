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

package sheet

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"unsafe"
)

type Regions struct {
	keys    map[string]*Region
	texture *common.Texture
	Count   int
}

func NewRegions() *Regions {
	return &Regions{
		keys: map[string]*Region{},
	}
}

func (s *Regions) SetTexture(texture *common.Texture) {
	s.Delete()
	s.texture = texture
}

func (s *Regions) Bind() {
	if s.texture != nil {
		s.texture.Bind()
	}
}

func (s *Regions) Unbind() {
	if s.texture != nil {
		s.texture.Unbind()
	}
}

func (s *Regions) Delete() {
	if s.texture != nil {
		s.texture.Delete()
		s.texture = nil
	}
}

func (s *Regions) AddRegion(key string, bounds, offset mgl32.Vec2) (out *Region) {
	var index int
	index = s.Count
	out = &Region{
		index:  index,
		bounds: bounds,
		offset: offset,
	}
	s.keys[key] = out
	s.Count++
	return
}

func (s *Regions) RegionExists(key string) (exists bool) {
	_, exists = s.keys[key]
	return
}

func (s *Regions) Region(key string) (out *Region, err error) {
	var exists bool
	if out, exists = s.keys[key]; !exists {
		err = fmt.Errorf("Invalid tile key %v", key)
		return
	}
	return
}

func (s *Regions) Upload(ubo *common.UniformBuffer) (err error) {
	var (
		region *Region
		entry  uniformRegion
		data   = make([]uniformRegion, s.Count)
		size   = s.Count * int(unsafe.Sizeof(entry))
	)
	if s.texture == nil {
		err = fmt.Errorf("No texture associated with sheet")
		return
	}
	for _, region = range s.keys {
		data[region.index] = region.textureBounds(s.texture.Size)
	}
	ubo.Upload(data, size)
	return
}

func (s *Regions) Texture() *common.Texture {
	return s.texture
}
