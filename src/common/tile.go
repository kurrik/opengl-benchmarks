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

package common

import (
	"fmt"
	"image"
	"unsafe"
)

type TileData [8]float32

func NewTileData(texW, texH, texX, texY, pxW, pxH, pxX, pxY float32) TileData {
	return TileData{texW, texH, texX, texY, pxW, pxH, pxX, pxY}
}

func (fd TileData) ImageBounds() image.Rectangle {
	var (
		x = int(fd[6])
		y = int(fd[7])
		w = int(fd[4])
		h = int(fd[5])
	)
	return image.Rectangle{
		image.Point{x, y},
		image.Point{x + w, y + h},
	}
}

type TileUniform struct {
	Tiles []TileData
	Count int
}

func NewTileUniform() *TileUniform {
	return &TileUniform{
		Tiles: []TileData{},
		Count: 0,
	}
}

func (td *TileUniform) AppendTile(tile TileData) (index int) {
	index = td.Count
	td.Count++
	td.Tiles = append(td.Tiles, tile)
	return
}

func (td *TileUniform) TileBounds(index int) (out image.Rectangle, err error) {
	if index < 0 || index > td.Count {
		err = fmt.Errorf("Invalid frame %v", index)
		return
	}
	out = td.Tiles[index].ImageBounds()
	return
}

func (td *TileUniform) TileBytes() int {
	var (
		point TileData
	)
	return td.Count * int(unsafe.Sizeof(point))
}
