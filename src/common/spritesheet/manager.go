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
	"github.com/kurrik/opengl-benchmarks/common/render"
	"github.com/kurrik/opengl-benchmarks/common/sprites"
	"github.com/kurrik/opengl-benchmarks/common/tile"
)

type Config struct {
	JsonPath      string
	PixelsPerUnit float32
	MaxInstances  uint32
	Renderer      *render.Renderer
}

type Manager struct {
	*tile.Manager
	cfg   Config
	sheet *sprites.Sheet
}

func NewManager(cfg Config) (mgr *Manager, err error) {
	var (
		tileManager *tile.Manager
		sheet       *sprites.Sheet
		tp          TexturePacker
	)
	if tileManager, err = tile.NewManager(cfg.Renderer); err != nil {
		return
	}
	if sheet, err = tp.LoadJSONArray(cfg.JsonPath); err != nil {
		return
	}
	mgr = &Manager{
		Manager: tileManager,
		cfg:     cfg,
		sheet:   sheet,
	}
	return
}

func (m *Manager) SetFrame(instance *render.Instance, frame string) (err error) {
	var s *sprites.Sprite
	if instance == nil {
		return // No error
	}
	if s, err = m.sheet.Sprite(frame); err != nil {
		return
	}
	instance.Frame = s.Index()
	instance.SetScale(s.WorldDimensions(m.cfg.PixelsPerUnit).Vec3(1.0))
	instance.MarkChanged()
	instance.Key = frame
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

func (m *Manager) Render(camera *common.Camera) {
	m.Manager.Render(camera, m.sheet)
}

// TODO: Refactor so that sheet can be shared between multiple renderers / managers.
func (m *Manager) Regions() *sprites.Sheet {
	return m.sheet
}
