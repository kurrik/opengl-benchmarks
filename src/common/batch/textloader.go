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

package batch

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"strings"
)

type TextMapping struct {
	defaultTile int
	mapping     map[rune]int
	sheet       *tile.Sheet
}

func NewTextMapping(sheet *tile.Sheet, defaultTile string) (out *TextMapping, err error) {
	var index int
	out = &TextMapping{
		sheet:   sheet,
		mapping: map[rune]int{},
	}
	if index, err = sheet.TileIndex(defaultTile); err != nil {
		return
	}
	out.defaultTile = index
	return

}

func (m *TextMapping) Set(r rune, tile string) (err error) {
	var index int
	if index, err = m.sheet.TileIndex(tile); err != nil {
		return
	}
	m.mapping[r] = index
	return
}

func (m *TextMapping) Get(r rune) (index int) {
	var exists bool
	if index, exists = m.mapping[r]; !exists {
		index = m.defaultTile
	}
	return
}

type TextLoader struct {
	mapping *TextMapping
}

func NewTextLoader(mapping *TextMapping) *TextLoader {
	return &TextLoader{
		mapping: mapping,
	}
}

func (l *TextLoader) addTile(x, y float32, index int, batch *Batch) {
	var findex = float32(index)
	batch.Points = append(batch.Points, []batchPoint{
		batchPoint{
			Position: mgl32.Vec3{x, y, 0},
			Texture:  mgl32.Vec2{0, 0},
			Tile:     findex,
		},
		batchPoint{
			Position: mgl32.Vec3{x + 1, y + 1, 0},
			Texture:  mgl32.Vec2{1, 1},
			Tile:     findex,
		},
		batchPoint{
			Position: mgl32.Vec3{x, y + 1, 0},
			Texture:  mgl32.Vec2{0, 1},
			Tile:     findex,
		},
		batchPoint{
			Position: mgl32.Vec3{x, y, 0},
			Texture:  mgl32.Vec2{0, 0},
			Tile:     findex,
		},
		batchPoint{
			Position: mgl32.Vec3{x + 1, y, 0},
			Texture:  mgl32.Vec2{1, 0},
			Tile:     findex,
		},
		batchPoint{
			Position: mgl32.Vec3{x + 1, y + 1, 0},
			Texture:  mgl32.Vec2{1, 1},
			Tile:     findex,
		},
	}...)
}

func (l *TextLoader) Load(grid string) (out *Batch, err error) {
	var (
		lines []string
		line  string
		char  rune
		y     int
		x     int
	)
	lines = strings.Split(strings.TrimSpace(grid), "\n")
	if len(lines) == 0 {
		err = fmt.Errorf("No lines in input data")
		return
	}
	out = NewBatch(len(lines) * len(lines[0]) * 6)
	for y, line = range lines {
		for x, char = range line {
			l.addTile(float32(x), float32(y), l.mapping.Get(char), out)
		}
	}
	return
}
