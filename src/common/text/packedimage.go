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
	"image"
	"image/draw"
)

type PackedImage struct {
	Width     int
	Height    int
	img       draw.Image
	shelves   []*shelf
	locations map[string]int
	tiles     map[string]int
	count     int
	Data      []float32
}

func NewPackedImage(w, h int) (i *PackedImage) {
	return &PackedImage{
		Width:     w,
		Height:    h,
		img:       image.NewRGBA(image.Rect(0, 0, w, h)),
		shelves:   []*shelf{newShelf()},
		locations: map[string]int{},
		Data:      []float32{},
		count:     0,
		tiles:     map[string]int{},
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
	)
	if bounds, err = src.Bounds(key); err != nil {
		return
	}
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
		exists    bool
	)
	if _, exists = i.tiles[key]; exists {
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
	i.locations[key] = len(i.Data)
	i.tiles[key] = i.count
	i.count += 1
	i.Data = append(i.Data,
		float32(w)/float32(texBounds.Max.X),
		float32(h)/float32(texBounds.Max.Y),
		float32(x)/float32(texBounds.Max.X),
		1.0-float32(y+h)/float32(texBounds.Max.Y),
		float32(w),
		float32(h),
		float32(destRect.Min.X),
		float32(destRect.Min.Y),
	)
	if glog.V(2) {
		glog.Infof("packRegion(%v): dest %v src %v", key, destRect, srcBounds.Min)
	}
	draw.Draw(i.img, destRect, src, srcBounds.Min, draw.Src)
	return
}

func (i *PackedImage) Bounds(key string) (out image.Rectangle, err error) {
	var (
		index int
		ok    bool
	)
	if index, ok = i.locations[key]; !ok {
		err = fmt.Errorf("Packed image did not contain key %v", key)
		return
	}
	out = image.Rectangle{
		image.Point{int(i.Data[index+6]), int(i.Data[index+7])},
		image.Point{int(i.Data[index+6] + i.Data[index+4]), int(i.Data[index+7] + i.Data[index+5])},
	}
	return
}

func (i *PackedImage) Tile(key string) (out int, err error) {
	var (
		ok bool
	)
	if out, ok = i.tiles[key]; !ok {
		err = fmt.Errorf("Packed image did not contain key %v", key)
		return
	}
	return
}
