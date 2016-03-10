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

package text

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"image"
	"image/draw"
)

type PackedImage struct {
	Width     int
	Height    int
	img       draw.Image
	shelves   []*shelf
	Sheet     *tile.Sheet
	pxPerUnit float32
}

func NewPackedImage(w, h int, pxPerUnit float32) (i *PackedImage) {
	return &PackedImage{
		Width:     w,
		Height:    h,
		img:       image.NewRGBA(image.Rect(0, 0, w, h)),
		shelves:   []*shelf{newShelf()},
		Sheet:     tile.NewSheet(w, h),
		pxPerUnit: pxPerUnit,
	}
}

func (i *PackedImage) Image() image.Image {
	return i.img
}

func (i *PackedImage) Pack(key string, img image.Image) (err error) {
	return i.packRegion(key, img, img.Bounds())
}

func (i *PackedImage) Copy(key string, src *PackedImage) (err error) {
	var (
		bounds image.Rectangle
		t      *tile.Tile
	)
	if t, err = src.Sheet.Tile(key); err != nil {
		return
	}
	bounds = t.ImageBounds()
	return i.packRegion(key, src.img, bounds)
}

func (i *PackedImage) packRegion(key string, src image.Image, srcBounds image.Rectangle) (err error) {
	var (
		j         int
		s         *shelf
		score     int
		bestScore int             = -1
		bestShelf int             = -1
		texBounds image.Rectangle = i.img.Bounds()
		w         int             = srcBounds.Max.X - srcBounds.Min.X
		h         int             = srcBounds.Max.Y - srcBounds.Min.Y
		maxW      int             = texBounds.Max.X
	)
	if i.Sheet.TileExists(key) {
		// Don't need to pack since it's already in here
		return
	}
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
		s = i.shelves[len(i.shelves)-1]
		if s.y+s.height+h > texBounds.Max.Y {
			// New shelf would exceed current image size
			err = fmt.Errorf("Cannot fit text into texture")
			return
		}
		i.shelves = append(i.shelves, i.shelves[len(i.shelves)-1].Close())
		bestShelf = len(i.shelves) - 1
	}
	s = i.shelves[bestShelf]
	var (
		x, y     = s.Add(w, h)
		destPt   = image.Pt(x, y)
		destRect = image.Rectangle{destPt, destPt.Add(image.Pt(w, h))}
	)
	i.Sheet.AddTile(key, w, h, x, y)
	/*
		tile.NewTile(
			float32(w)/float32(texBounds.Max.X),
			float32(h)/float32(texBounds.Max.Y),
			float32(x)/float32(texBounds.Max.X),
			1.0-float32(y+h)/float32(texBounds.Max.Y),
			float32(w)/i.pxPerUnit,
			float32(h)/i.pxPerUnit,
			float32(destRect.Min.X),
			float32(destRect.Min.Y),
		))
	*/
	if glog.V(2) {
		glog.Infof("packRegion(%v): dest %v src %v", key, destRect, srcBounds.Min)
	}
	draw.Draw(i.img, destRect, src, srcBounds.Min, draw.Src)
	return
}
