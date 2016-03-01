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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/binpacking"
	"image/draw"
)

type Manager struct {
	PackedImage   *binpacking.PackedImage
	nextID        TextID
	packedTexture *common.Texture
	instances     map[TextID]*managerInstance
	rendererData  []rendererInstance
	renderer      *Renderer
	maxInstances  uint32
}

func NewManager(maxInstances uint32) (mgr *Manager, err error) {
	mgr = &Manager{
		PackedImage:  binpacking.NewPackedImage(512, 512),
		instances:    map[TextID]*managerInstance{},
		maxInstances: maxInstances,
		rendererData: make([]rendererInstance, maxInstances),
	}
	if mgr.renderer, err = NewRenderer(); err != nil {
		return
	}
	return
}

func (m *Manager) CreateText() (id TextID) {
	id = m.nextID
	m.nextID += 1
	return
}

func (m *Manager) ensureInstance(id TextID) {
	var exists bool
	if _, exists = m.instances[id]; !exists {
		m.instances[id] = &managerInstance{}
	}
}

func (m *Manager) SetText(id TextID, text string, font *FontFace) (err error) {
	var (
		img draw.Image
		key = fmt.Sprintf("%v", id)
	)
	if img, err = font.GetImage(text); err != nil {
		return
	}
	m.ensureInstance(id)
	m.PackedImage.Pack(key, img)
	if m.instances[id].packedIndex, err = m.PackedImage.Index(key); err != nil {
		return
	}
	if m.packedTexture, err = common.GetTexture(
		m.PackedImage.Image(),
		common.SmoothingLinear,
	); err != nil {
		return
	}
	return
}

func (m *Manager) SetPosition(id TextID, position mgl32.Vec2) {
	m.ensureInstance(id)
	m.instances[id].position = position
}

func (m *Manager) Bind() {
	if m.packedTexture != nil {
		m.packedTexture.Bind()
	}
	m.renderer.Bind()
}

func (m *Manager) Unbind() {
	if m.packedTexture != nil {
		m.packedTexture.Unbind()
	}
	m.renderer.Unbind()
}

func (m *Manager) Delete() {
	if m.packedTexture != nil {
		m.packedTexture.Delete()
		m.packedTexture = nil
	}
	m.renderer.Delete()
}

func (m *Manager) Render(camera *common.Camera) {
	// Temporary:
	var (
		scale  = mgl32.Scale3D(1.0/128.0, 1.0/128.0, 1.0)
		rot1   = mgl32.HomogRotate3DZ(mgl32.DegToRad(5.0))
		trans2 = mgl32.Translate3D(1, 1, 0)
		rot2   = mgl32.HomogRotate3DZ(mgl32.DegToRad(15.0))
	)
	data := &rendererData{
		Instances: []rendererInstance{
			rendererInstance{
				model: rot1.Mul4(scale),
				tile:  2,
			},
			rendererInstance{
				model: trans2.Mul4(rot2).Mul4(scale),
				tile:  7,
			},
		},
	}
	m.renderer.Render(camera, data, m.PackedImage.Data)
}
