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
	"image"
)

type Tile [8]float32

func NewTile(texW, texH, texX, texY, pxW, pxH, pxX, pxY float32) Tile {
	return Tile{texW, texH, texX, texY, pxW, pxH, pxX, pxY}
}

func (t Tile) ImageBounds() image.Rectangle {
	var (
		x = int(t[6])
		y = int(t[7])
		w = int(t[4])
		h = int(t[5])
	)
	return image.Rectangle{
		image.Point{x, y},
		image.Point{x + w, y + h},
	}
}
