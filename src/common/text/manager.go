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

type Config struct {
	MaxInstances  uint32
	TextureWidth  int
	TextureHeight int
}

type Manager struct {
	*tile.Manager
	cfg         Config
	PackedImage *PackedImage
	texture     *common.Texture
}

func NewManager(cfg Config) (mgr *Manager, err error) {
	var (
		tileManager *tile.Manager
	)
	if tileManager, err = tile.NewManager(tile.Config{
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

func (m *Manager) SetText(instance *tile.Instance, text string, font *FontFace) (err error) {
	var (
		img draw.Image
	)
	if instance == nil {
		return // No error.
	}
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
	if m.texture != nil {
		m.texture.Delete()
	}
	if m.texture, err = common.GetTexture(
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
		instance *tile.Instance
	)
	if glog.V(1) {
		glog.Info("Repacking image")
	}
	newImage = NewPackedImage(m.PackedImage.Width, m.PackedImage.Height)
	instance = m.Instances.Head()
	for instance != nil {
		if err = newImage.Copy(instance.Key, m.PackedImage); err != nil {
			return
		}
		if instance.Tile, err = newImage.Sheet.TileIndex(instance.Key); err != nil {
			return
		}
		instance.Dirty = true
		instance = instance.Next()
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
	if m.texture != nil {
		m.texture.Bind()
	}
	m.Manager.Bind()
}

func (m *Manager) Unbind() {
	if m.texture != nil {
		m.texture.Unbind()
	}
	m.Manager.Unbind()
}

func (m *Manager) Delete() {
	if m.texture != nil {
		m.texture.Delete()
		m.texture = nil
	}
	m.Manager.Delete()
}

func (m *Manager) Render(camera *common.Camera) {
	m.Manager.Render(camera, m.PackedImage.Sheet)
}
