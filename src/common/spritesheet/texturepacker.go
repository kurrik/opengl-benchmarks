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

package spritesheet

import (
	"encoding/json"
	"github.com/kurrik/opengl-benchmarks/common/tile"
)

type texturePackerFloatCoords struct {
	X float32 `json:x,omitempty`
	Y float32 `json:y,omitempty`
}

type texturePackerIntCoords struct {
	X int `json:x,omitempty`
	Y int `json:y,omitempty`
	W int `json:w,omitempty`
	H int `json:h,omitempty`
}

type texturePackerFrame struct {
	Filename         string                   `json:filename`
	Frame            texturePackerIntCoords   `json:frame`
	Rotated          bool                     `json:rotated`
	Trimmed          bool                     `json:trimmed`
	SpriteSourceSize texturePackerIntCoords   `json:spriteSourceSize`
	SourceSize       texturePackerIntCoords   `json:sourceSize`
	Pivot            texturePackerFloatCoords `json:pivot`
}

type texturePackerMeta struct {
	Image  string                 `json:image`
	Format string                 `json:format`
	Size   texturePackerIntCoords `json:size`
	Scale  string                 `json:scale`
}

type texturePackerJSONArray struct {
	Frames []texturePackerFrame `json:frames`
	Meta   texturePackerMeta    `json:meta`
}

func (f texturePackerFrame) ToTile(meta texturePackerMeta, pxPerUnit float32) tile.Tile {
	var (
		//sourceW          = float32(f.SpriteSourceSize.W)
		//sourceH          = float32(f.SpriteSourceSize.H)
		//sourceX          = float32(f.SpriteSourceSize.X)
		//sourceY          = float32(f.SpriteSourceSize.Y)
		textureX         = float32(f.Frame.X)
		textureY         = float32(f.Frame.Y)
		textureW         = float32(f.Frame.W)
		textureH         = float32(f.Frame.H)
		textureOriginalW = float32(meta.Size.W)
		textureOriginalH = float32(meta.Size.H)
		texX             = textureX / textureOriginalW
		texY             = textureY / textureOriginalH
		texW             = textureW / textureOriginalW
		texH             = textureH / textureOriginalH
		//ptW              = sourceW / pxPerUnit
		//ptH              = sourceH / pxPerUnit
	)
	return tile.NewTile(
		texW,
		texH,
		texX,
		texY,
		textureW,
		textureH,
		textureX,
		textureY,
	)
	/*
		return &SpritesheetFrame{
			Size:          mgl32.Vec2{ptW, ptH},
			TextureOffset: mgl32.Vec2{texX, texY},
			TextureSize:   mgl32.Vec2{texW, texH},
		}
	*/
}

func ParseTexturePackerJSONArrayString(contents string, pxPerUnit float32) (s *tile.Sheet, texturePath string, err error) {
	var (
		parsed texturePackerJSONArray
	)
	if err = json.Unmarshal([]byte(contents), &parsed); err != nil {
		return
	}
	texturePath = parsed.Meta.Image
	s = tile.NewSheet()
	for _, frame := range parsed.Frames {
		s.AddTile(
			frame.Filename,
			frame.ToTile(parsed.Meta, pxPerUnit),
		)
	}
	return
}
