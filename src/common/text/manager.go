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
	"github.com/kurrik/opengl-benchmarks/common/sheet"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"image/draw"
)

type Config struct {
	MaxInstances  uint32
	TextureWidth  int
	TextureHeight int
	PixelsPerUnit float32
}

type Manager struct {
	*tile.Manager
	cfg     Config
	regions *sheet.PackedRegions
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
		cfg:     cfg,
		Manager: tileManager,
		regions: sheet.NewPackedRegions(
			cfg.TextureWidth,
			cfg.TextureHeight,
		),
	}
	return
}

func (m *Manager) SetText(instance *tile.Instance, text string, font *FontFace) (err error) {
	var (
		img    draw.Image
		region *sheet.Region
	)
	if instance == nil {
		return // No error.
	}
	if img, err = font.GetImage(text); err != nil {
		return
	}
	if err = m.regions.Pack(text, img); err != nil {
		// Attempt to compact the texture.
		if err = m.repackImage(); err != nil {
			return
		}
		if err = m.regions.Pack(text, img); err != nil {
			return
		}
	}
	if region, err = m.regions.Regions.Region(text); err != nil {
		return
	}
	instance.Tile = region.Index()
	instance.SetScale(region.WorldDimensions(m.cfg.PixelsPerUnit).Vec3(1.0))
	instance.Dirty = true
	instance.Key = text
	if err = m.generateTexture(); err != nil {
		return
	}
	return
}

func (m *Manager) generateTexture() (err error) {
	var (
		texture *common.Texture
	)
	if texture, err = common.GetTexture(
		m.regions.Image(),
		common.SmoothingLinear,
	); err != nil {
		return
	}
	m.regions.SetTexture(texture)
	return
}

func (m *Manager) repackImage() (err error) {
	var (
		newImage *sheet.PackedRegions
		instance *tile.Instance
		region   *sheet.Region
	)
	if glog.V(1) {
		glog.Info("Repacking image")
	}
	newImage = sheet.NewPackedRegions(
		m.regions.Width,
		m.regions.Height,
	)
	instance = m.Instances.Head()
	for instance != nil {
		if err = newImage.Copy(instance.Key, m.regions); err != nil {
			return
		}
		if region, err = newImage.Regions.Region(instance.Key); err != nil {
			return
		}
		instance.Tile = region.Index()
		instance.Dirty = true
		instance = instance.Next()
	}
	m.regions = newImage
	if err = m.generateTexture(); err != nil {
		return
	}
	if glog.V(1) {
		glog.Info("Done repacking")
	}
	return
}

func (m *Manager) Bind() {
	m.regions.Bind()
	m.Manager.Bind()
}

func (m *Manager) Unbind() {
	m.regions.Unbind()
	m.Manager.Unbind()
}

func (m *Manager) Delete() {
	m.regions.Delete()
	m.Manager.Delete()
	m.regions = nil
	m.Manager = nil
}

func (m *Manager) Render(camera *common.Camera) {
	m.Manager.Render(camera, m.regions)
}

func (m *Manager) Regions() *sheet.PackedRegions {
	return m.regions
}
