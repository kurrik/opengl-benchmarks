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
	"fmt"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/binpacking"
	"image/draw"
)

type Manager struct {
	nextID        TextID
	packedImage   *binpacking.PackedImage
	packedTexture *common.Texture
	indices       map[TextID]int
}

func NewManager() *Manager {
	return &Manager{
		packedImage: binpacking.NewPackedImage(512, 512),
		indices:     map[TextID]int{},
	}
}

func (t *Manager) CreateText() (id TextID) {
	id = t.nextID
	t.nextID += 1
	return
}

func (t *Manager) SetText(id TextID, text string, font *FontFace) (err error) {
	var (
		img draw.Image
		key = fmt.Sprintf("%v", id)
	)
	if img, err = font.GetImage(text); err != nil {
		return
	}
	t.packedImage.Pack(key, img)
	if t.indices[id], err = t.packedImage.Index(key); err != nil {
		return
	}
	if t.packedTexture, err = common.GetTexture(
		t.packedImage.Image(),
		common.SmoothingLinear,
	); err != nil {
		return
	}
	return
}

func (t *Manager) Bind() {
	if t.packedTexture != nil {
		t.packedTexture.Bind()
	}
}

func (t *Manager) Unbind() {
	if t.packedTexture != nil {
		t.packedTexture.Unbind()
	}
}

func (t *Manager) Delete() {
	if t.packedTexture != nil {
		t.packedTexture.Delete()
		t.packedTexture = nil
	}
}
