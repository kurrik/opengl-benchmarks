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

package text

import (
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/binpacking"
	"image/draw"
)

type Text struct {
	packedImage   *binpacking.PackedImage
	packedTexture *common.Texture
}

func NewText() *Text {
	return &Text{
		packedImage: binpacking.NewPackedImage(512, 512),
	}
}

func (t *Text) Set(key, text string, font *FontFace) (err error) {
	var img draw.Image
	if img, err = font.GetImage(text); err != nil {
		return
	}
	t.packedImage.Pack(key, img)
	if t.packedTexture, err = common.GetTexture(
		t.packedImage.Image(),
		common.SmoothingLinear,
	); err != nil {
		return
	}
	return
}

func (t *Text) Bind() {
	if t.packedTexture != nil {
		t.packedTexture.Bind()
	}
}

func (t *Text) Unbind() {
	if t.packedTexture != nil {
		t.packedTexture.Unbind()
	}
}

func (t *Text) Delete() {
	if t.packedTexture != nil {
		t.packedTexture.Delete()
		t.packedTexture = nil
	}
}
