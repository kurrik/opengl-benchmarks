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
	"github.com/golang/glog"
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"image/draw"
)

type ManagerConfig struct {
	MaxInstances  uint32
	TextureWidth  int
	TextureHeight int
}

type Manager struct {
	*tile.Manager
	cfg           ManagerConfig
	PackedImage   *PackedImage
	packedTexture *common.Texture
}

func NewManager(cfg ManagerConfig) (mgr *Manager, err error) {
	var (
		tileManager *tile.Manager
	)
	if tileManager, err = tile.NewManager(tile.ManagerConfig{
		MaxInstances: cfg.MaxInstances,
	}); err != nil {
		return
	}
	mgr = &Manager{
		Manager:     tileManager,
		PackedImage: NewPackedImage(cfg.TextureWidth, cfg.TextureHeight),
		cfg:         cfg,
	}
	return
}

func (m *Manager) SetText(id tile.InstanceID, text string, font *FontFace) (err error) {
	var (
		img      draw.Image
		instance *tile.TileInstance
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
	if instance, err = m.GetInstance(id); err != nil {
		return
	}
	if instance.Tile, err = m.PackedImage.Sheet.TileIndex(text); err != nil {
		return
	}
	instance.Dirty = true
	instance.Key = text
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
		instance *tile.TileInstance
	)
	if glog.V(1) {
		glog.Info("Repacking image")
	}
	newImage = NewPackedImage(m.PackedImage.Width, m.PackedImage.Height)
	for _, instance = range m.Instances {
		if err = newImage.Copy(instance.Key, m.PackedImage); err != nil {
			return
		}
		if instance.Tile, err = newImage.Sheet.TileIndex(instance.Key); err != nil {
			return
		}
		instance.Dirty = true
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

func (m *Manager) Bind() {
	if m.packedTexture != nil {
		m.packedTexture.Bind()
	}
	m.Manager.Bind()
}

func (m *Manager) Unbind() {
	if m.packedTexture != nil {
		m.packedTexture.Unbind()
	}
	m.Manager.Unbind()
}

func (m *Manager) Delete() {
	if m.packedTexture != nil {
		m.packedTexture.Delete()
		m.packedTexture = nil
	}
	m.Manager.Delete()
}

func (m *Manager) Render(camera *common.Camera) {
	m.Manager.Render(camera, m.PackedImage.Sheet)
}
