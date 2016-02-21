// Copyright 2015 Arne Roomann-Kurrik
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
	"image"
	"image/draw"
)

type shelf struct {
	x      int
	y      int
	height int
	isOpen bool
}

func newShelf() *shelf {
	return &shelf{
		x:      0,
		y:      0,
		height: 0,
		isOpen: true,
	}
}

func (s *shelf) FitsX(w, maxW int) bool {
	return s.x+w <= maxW
}

func (s *shelf) FitsY(h int) bool {
	return s.height >= h
}

func (s *shelf) RemainingX(maxW int) int {
	return maxW - s.x
}

func (s *shelf) CanAdd(w, h, maxW int) bool {
	if !s.FitsX(w, maxW) {
		return false
	}
	if !s.isOpen && !s.FitsY(h) {
		return false
	}
	return true
}

func (s *shelf) Add(w, h int) (origX, origY int) {
	origX = s.x
	origY = s.y
	if s.height < h {
		s.height = h
	}
	s.x += w
	return
}

func (s *shelf) Close() (out *shelf) {
	out = newShelf()
	out.y = s.y + s.height
	s.isOpen = false
	return
}

func (s *shelf) BestAreaFit(w, h, maxW int) int {
	var (
		shelfArea = s.RemainingX(maxW) * s.height
		wordArea  = w * h
	)
	return shelfArea - wordArea
}

type ImagePacked struct {
	img     draw.Image
	shelves []*shelf
}

func NewImagePacked(w, h int) (i *ImagePacked) {
	return &ImagePacked{
		img:     image.NewRGBA(image.Rect(0, 0, w, h)),
		shelves: []*shelf{newShelf()},
	}
}

func (i *ImagePacked) Image() image.Image {
	return i.img
}

func (i *ImagePacked) Pack(img image.Image) {
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
	draw.Draw(i.img, destRect, img, imgBounds.Min, draw.Src)
}
