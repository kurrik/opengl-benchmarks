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
	"github.com/go-gl/mathgl/mgl32"
	"image"
)

type Region struct {
	index  int
	bounds mgl32.Vec2
	offset mgl32.Vec2
}

func (r *Region) Index() int {
	return r.index
}

func (r *Region) ImageBounds() image.Rectangle {
	return image.Rectangle{
		image.Point{int(r.offset.X()), int(r.offset.Y())},
		image.Point{int(r.offset.X() + r.bounds.X()), int(r.offset.Y() + r.bounds.Y())},
	}
}

func (r *Region) textureBounds(textureBounds mgl32.Vec2) uniformRegion {
	return newUniformRegion(
		r.bounds.X()/textureBounds.X(),
		r.bounds.Y()/textureBounds.Y(),
		r.offset.X()/textureBounds.X(),
		1.0-(r.offset.Y()+r.bounds.Y()-1.0)/textureBounds.Y(),
	)
}

func (r *Region) WorldDimensions(pxPerUnit float32) mgl32.Vec2 {
	return mgl32.Vec2{
		r.bounds.X() / pxPerUnit,
		r.bounds.Y() / pxPerUnit,
	}
}
