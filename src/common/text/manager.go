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
	"github.com/kurrik/opengl-benchmarks/common/render"
	"github.com/kurrik/opengl-benchmarks/common/sprites"
	"image/draw"
)

type Config struct {
	TextureWidth  int
	TextureHeight int
	PixelsPerUnit float32
}

type Manager struct {
	Instances *render.InstanceList
	cfg       Config
	sheet     *sprites.PackedSheet
}

func NewManager(cfg Config) *Manager {
	return &Manager{
		Instances: render.NewInstanceList(),
		cfg:       cfg,
		sheet: sprites.NewPackedSheet(
			cfg.TextureWidth,
			cfg.TextureHeight,
		),
	}
}

func (m *Manager) SetText(instance *render.Instance, text string, font *FontFace) (err error) {
	var (
		img    draw.Image
		sprite *sprites.Sprite
	)
	if instance == nil {
		return // No error.
	}
	if img, err = font.GetImage(text); err != nil {
		return
	}
	if err = m.sheet.Pack(text, img); err != nil {
		// Attempt to compact the texture.
		if err = m.repackImage(); err != nil {
			return
		}
		if err = m.sheet.Pack(text, img); err != nil {
			return
		}
	}
	if sprite, err = m.sheet.Sprite(text); err != nil {
		return
	}
	instance.Frame = sprite.Index()
	instance.SetScale(sprite.WorldDimensions(m.cfg.PixelsPerUnit).Vec3(1.0))
	instance.MarkChanged()
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
		m.sheet.Image(),
		common.SmoothingLinear,
	); err != nil {
		return
	}
	m.sheet.SetTexture(texture)
	return
}

func (m *Manager) repackImage() (err error) {
	var (
		newImage *sprites.PackedSheet
		instance *render.Instance
		sprite   *sprites.Sprite
	)
	if glog.V(1) {
		glog.Info("Repacking image")
	}
	newImage = sprites.NewPackedSheet(
		m.sheet.Width,
		m.sheet.Height,
	)
	instance = m.Instances.Head()
	for instance != nil {
		if err = newImage.Copy(instance.Key, m.sheet); err != nil {
			return
		}
		if sprite, err = newImage.Sheet.Sprite(instance.Key); err != nil {
			return
		}
		instance.Frame = sprite.Index()
		instance.MarkChanged()
		instance = instance.Next()
	}
	m.sheet = newImage
	if err = m.generateTexture(); err != nil {
		return
	}
	if glog.V(1) {
		glog.Info("Done repacking")
	}
	return
}

func (m *Manager) Bind() {
	m.sheet.Bind()
}

func (m *Manager) Unbind() {
	m.sheet.Unbind()
}

func (m *Manager) Delete() {
	m.sheet.Delete()
	m.sheet = nil
}

func (m *Manager) Regions() *sprites.PackedSheet {
	return m.sheet
}
