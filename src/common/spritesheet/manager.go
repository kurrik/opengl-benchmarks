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
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/tile"
	"io/ioutil"
	"path"
)

type Config struct {
	JsonPath      string
	PixelsPerUnit float32
	MaxInstances  uint32
}

type Manager struct {
	*tile.Manager
	cfg     Config
	sheet   *tile.Sheet
	texture *common.Texture
}

func NewManager(cfg Config) (mgr *Manager, err error) {
	var (
		data        []byte
		dir         = path.Dir(cfg.JsonPath)
		texture     *common.Texture
		texturePath string
		tileManager *tile.Manager
		sheet       *tile.Sheet
	)
	if tileManager, err = tile.NewManager(tile.Config{
		MaxInstances: cfg.MaxInstances,
	}); err != nil {
		return
	}
	if data, err = ioutil.ReadFile(cfg.JsonPath); err != nil {
		return
	}
	if sheet, texturePath, err = ParseTexturePackerJSONArrayString(
		string(data),
		cfg.PixelsPerUnit,
	); err != nil {
		return
	}
	if texture, err = common.LoadTexture(
		path.Join(dir, texturePath),
		common.SmoothingNearest,
	); err != nil {
		return
	}
	mgr = &Manager{
		Manager: tileManager,
		sheet:   sheet,
		texture: texture,
		cfg:     cfg,
	}
	return
}

func (m *Manager) SetFrame(instance *tile.Instance, frame string) (err error) {
	var t *tile.Tile
	if instance == nil {
		return // No error
	}
	if t, err = m.sheet.Tile(frame); err != nil {
		return
	}
	instance.Tile = t.Index()
	instance.SetScale(t.WorldDimensions(m.cfg.PixelsPerUnit).Vec3(1.0))
	instance.Dirty = true
	instance.Key = frame
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
	m.Manager.Render(camera, m.sheet)
}

// TODO: Refactor so that sheet can be shared between multiple renderers / managers.
func (m *Manager) GetSheet() *tile.Sheet {
	return m.sheet
}

// TODO: Refactor so that texture can be shared between multiple renderers / managers.
func (m *Manager) GetTexture() *common.Texture {
	return m.texture
}
