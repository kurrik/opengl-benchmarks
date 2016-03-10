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
	"github.com/go-gl/mathgl/mgl32"
	"image"
)

type Tile struct {
	index int
	pxW   int
	pxH   int
	pxX   int
	pxY   int
}

func (t *Tile) Index() int {
	return t.index
}

func (t *Tile) ImageBounds() image.Rectangle {
	return image.Rectangle{
		image.Point{t.pxX, t.pxY},
		image.Point{t.pxX + t.pxW, t.pxY + t.pxH},
	}
}

func (t *Tile) textureBounds(texW, texH float32) rTile {
	return newRTile(
		float32(t.pxW) / texW,
		float32(t.pxH) / texH,
		float32(t.pxX) / texW,
		1.0 - float32(t.pxY+t.pxH-1)/texH,
	)
}

func (t *Tile) WorldDimensions(pxPerUnit float32) mgl32.Vec2 {
	return mgl32.Vec2{
		float32(t.pxW) / pxPerUnit,
		float32(t.pxH) / pxPerUnit,
	}
}
