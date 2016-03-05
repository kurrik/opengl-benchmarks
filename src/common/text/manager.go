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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/glog"
	"github.com/kurrik/opengl-benchmarks/common"
	"image/draw"
)

type ManagerConfig struct {
	MaxInstances  uint32
	TextureWidth  int
	TextureHeight int
}

type Manager struct {
	cfg           ManagerConfig
	PackedImage   *PackedImage
	nextID        ID
	packedTexture *common.Texture
	renderer      *Renderer
	instances     map[ID]*Instance
	rendererData  rendererData
}

func NewManager(cfg ManagerConfig) (mgr *Manager, err error) {
	mgr = &Manager{
		PackedImage: NewPackedImage(cfg.TextureWidth, cfg.TextureHeight),
		cfg:         cfg,
		instances:   map[ID]*Instance{},
		rendererData: rendererData{
			Count:     0,
			Instances: make([]rendererInstance, cfg.MaxInstances),
		},
	}
	if mgr.renderer, err = NewRenderer(); err != nil {
		return
	}
	return
}

func (m *Manager) CreateText() (id ID, err error) {
	if uint32(m.rendererData.Count) >= m.cfg.MaxInstances {
		err = fmt.Errorf("Max text instances reached")
		return
	}
	id = m.nextID
	m.instances[id] = &Instance{
		renderIndex: m.rendererData.Count,
		tile:        0,
		position:    mgl32.Vec3{0, 0, 0},
		rotation:    0,
		dirty:       true,
		Text:        "",
	}
	m.nextID += 1
	m.rendererData.Count += 1
	return
}

func (m *Manager) getInstance(id ID) (instance *Instance, err error) {
	var (
		exists bool
	)
	if instance, exists = m.instances[id]; !exists {
		err = fmt.Errorf("Invalid text instance ID: %v", id)
		return
	}
	return
}

func (m *Manager) SetText(id ID, text string, font *FontFace) (err error) {
	var (
		img      draw.Image
		instance *Instance
	)
	if img, err = font.GetImage(text); err != nil {
		return
	}
	if err = m.PackedImage.Pack(text, img); err != nil {
		// Attempt to compact the texture.
		if err = m.repackImage(); err != nil {
			return
		}
		if err = m.PackedImage.Pack(text, img); err != nil {
			return
		}
	}
	if instance, err = m.getInstance(id); err != nil {
		return
	}
	if instance.tile, err = m.PackedImage.Sheet.TileIndex(text); err != nil {
		return
	}
	instance.dirty = true
	instance.Text = text
	if err = m.generateTexture(); err != nil {
		return
	}
	return
}

func (m *Manager) generateTexture() (err error) {
	if m.packedTexture != nil {
		m.packedTexture.Delete()
	}
	if m.packedTexture, err = common.GetTexture(
		m.PackedImage.Image(),
		common.SmoothingLinear,
	); err != nil {
		return
	}
	return
}

func (m *Manager) repackImage() (err error) {
	var (
		newImage *PackedImage
		instance *Instance
	)
	if glog.V(1) {
		glog.Info("Repacking image")
	}
	newImage = NewPackedImage(m.PackedImage.Width, m.PackedImage.Height)
	for _, instance = range m.instances {
		if err = newImage.Copy(instance.Text, m.PackedImage); err != nil {
			return
		}
		if instance.tile, err = newImage.Sheet.TileIndex(instance.Text); err != nil {
			return
		}
		instance.dirty = true
	}
	m.PackedImage = newImage
	if err = m.generateTexture(); err != nil {
		return
	}
	if glog.V(1) {
		glog.Info("Done repacking")
	}
	return
}

func (m *Manager) GetInstance(id ID) (instance *Instance, err error) {
	if instance, err = m.getInstance(id); err != nil {
		return
	}
	return
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
	var (
		instance  *Instance
		rInstance *rendererInstance
		scale     mgl32.Mat4
		rot       mgl32.Mat4
		trans     mgl32.Mat4
	)
	scale = mgl32.Scale3D(
		1.0/camera.PxPerUnit.X(),
		1.0/camera.PxPerUnit.Y(),
		1.0,
	)
	for _, instance = range m.instances {
		if instance.dirty {
			rInstance = &m.rendererData.Instances[instance.renderIndex]
			rInstance.tile = float32(instance.tile)
			rot = mgl32.HomogRotate3DZ(mgl32.DegToRad(instance.rotation))
			trans = mgl32.Translate3D(
				instance.position.X(),
				instance.position.Y(),
				instance.position.Z(),
			)
			rInstance.model = trans.Mul4(rot).Mul4(scale)
			instance.dirty = false
		}
	}
	m.renderer.Render(camera, &m.rendererData, m.PackedImage.Sheet)
}
