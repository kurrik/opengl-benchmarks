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

type Uniform struct {
	Tiles []Tile
	Count int
}

func NewUniform() *Uniform {
	return &Uniform{
		Tiles: []Tile{},
		Count: 0,
	}
}

func (u *Uniform) AppendTile(tile Tile) (index int) {
	index = u.Count
	u.Count++
	u.Tiles = append(u.Tiles, tile)
	return
}

func (u *Uniform) TileBounds(index int) (out image.Rectangle, err error) {
	if index < 0 || index > u.Count {
		err = fmt.Errorf("Invalid frame %v", index)
		return
	}
	out = u.Tiles[index].ImageBounds()
	return
}

func (u *Uniform) TileBytes() int {
	var (
		point Tile
	)
	return u.Count * int(unsafe.Sizeof(point))
}
