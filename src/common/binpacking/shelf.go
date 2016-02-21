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

import ()

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
