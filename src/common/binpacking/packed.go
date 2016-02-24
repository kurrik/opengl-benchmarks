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

package binpacking

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/draw"
)

type PackedImage struct {
	img       draw.Image
	shelves   []*shelf
	locations map[string]mgl32.Vec4
}

func NewPackedImage(w, h int) (i *PackedImage) {
	return &PackedImage{
		img:       image.NewRGBA(image.Rect(0, 0, w, h)),
		shelves:   []*shelf{newShelf()},
		locations: map[string]mgl32.Vec4{},
	}
}

func (i *PackedImage) Image() image.Image {
	return i.img
}

func (i *PackedImage) Pack(key string, img image.Image) {
	var (
		j         int
		s         *shelf
		score     int
		bestScore int             = -1
		bestShelf int             = -1
		imgBounds image.Rectangle = img.Bounds()
		texBounds image.Rectangle = i.img.Bounds()
		w         int             = imgBounds.Max.X
		h         int             = imgBounds.Max.Y
		maxW      int             = texBounds.Max.X
	)
	for j, s = range i.shelves {
		if s.CanAdd(w, h, maxW) {
			score = s.BestAreaFit(w, h, maxW)
			if score > bestScore {
				bestScore = score
				bestShelf = j
			}
		}
	}
	if bestShelf == -1 {
		i.shelves = append(i.shelves, i.shelves[len(i.shelves)-1].Close())
		bestShelf = len(i.shelves) - 1
	}
	s = i.shelves[bestShelf]
	var (
		x, y     = s.Add(w, h)
		destPt   = image.Pt(x, y)
		destRect = image.Rectangle{destPt, destPt.Add(imgBounds.Max)}
	)
	i.locations[key] = mgl32.Vec4{
		float32(destRect.Min.X),
		float32(destRect.Min.Y),
		float32(destRect.Max.X),
		float32(destRect.Max.Y),
	}
	draw.Draw(i.img, destRect, img, imgBounds.Min, draw.Src)
}

func (i *PackedImage) Bounds(key string) (out mgl32.Vec4, err error) {
	var (
		ok bool
	)
	if out, ok = i.locations[key]; !ok {
		err = fmt.Errorf("Packed image did not contain key %v", key)
		return
	}
	return
}
