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

package tile

import (
	"fmt"
	"github.com/kurrik/opengl-benchmarks/common"
	"unsafe"
)

type Sheet struct {
	keys    map[string]*Tile
	Count   int
	Version int
	width   int
	height  int
}

func NewSheet(w, h int) *Sheet {
	return &Sheet{
		keys:    map[string]*Tile{},
		width:   w,
		height:  h,
		Count:   0,
		Version: 0,
	}
}

func (s *Sheet) AddTile(key string, w, h, x, y int) (out *Tile) {
	var index int
	index = s.Count
	out = &Tile{
		index: index,
		pxW:   w,
		pxH:   h,
		pxX:   x,
		pxY:   y,
	}
	s.keys[key] = out
	s.Count++
	s.Version++
	return
}

func (s *Sheet) TileExists(key string) (exists bool) {
	_, exists = s.keys[key]
	return
}

func (s *Sheet) Tile(key string) (out *Tile, err error) {
	var exists bool
	if out, exists = s.keys[key]; !exists {
		err = fmt.Errorf("Invalid tile key %v", key)
		return
	}
	return
}

func (s *Sheet) Upload(ubo *common.UniformBuffer) {
	var (
		t     *Tile
		entry rTile
		data  = make([]rTile, s.Count)
		size  = s.Count * int(unsafe.Sizeof(entry))
	)
	for _, t = range s.keys {
		data[t.index] = t.textureBounds(float32(s.width), float32(s.height))
	}
	ubo.Upload(data, size)
}
