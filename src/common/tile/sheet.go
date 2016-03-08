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
	"image"
	"unsafe"
)

type Sheet struct {
	Tiles []Tile
	Count int
	keys  map[string]int
}

func NewSheet() *Sheet {
	return &Sheet{
		Tiles: []Tile{},
		Count: 0,
		keys:  map[string]int{},
	}
}

func (s *Sheet) AddTile(key string, tile Tile) {
	var index int
	index = s.Count
	s.Tiles = append(s.Tiles, tile)
	s.keys[key] = index
	s.Count++
}

func (s *Sheet) TileExists(key string) (exists bool) {
	_, exists = s.keys[key]
	return
}

func (s *Sheet) TileIndex(key string) (index int, err error) {
	var exists bool
	if index, exists = s.keys[key]; !exists {
		err = fmt.Errorf("Invalid frame %v", index)
		return
	}
	return
}

func (s *Sheet) TileBounds(key string) (out image.Rectangle, err error) {
	var index int
	if index, err = s.TileIndex(key); err != nil {
		return
	}
	out = s.Tiles[index].ImageBounds()
	return
}

func (s *Sheet) Bytes() int {
	var point Tile
	return s.Count * int(unsafe.Sizeof(point))
}
