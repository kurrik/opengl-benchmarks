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
	Sheet         *sprites.Sheet
	PixelsPerUnit float32
	MaxInstances  uint32
	Renderer      *render.Renderer
}

type Manager struct {
	*tile.Manager
	cfg Config
}

func NewManager(cfg Config) (mgr *Manager, err error) {
	var (
		tileManager *tile.Manager
	)
	if tileManager, err = tile.NewManager(cfg.Renderer); err != nil {
		return
	}
	mgr = &Manager{
		Manager: tileManager,
		cfg:     cfg,
	}
	return
}

func (m *Manager) SetFrame(instance *render.Instance, frame string) (err error) {
	var s *sprites.Sprite
	if instance == nil {
		return // No error
	}
	if s, err = m.cfg.Sheet.Sprite(frame); err != nil {
		return
	}
	instance.Frame = s.Index()
	instance.SetScale(s.WorldDimensions(m.cfg.PixelsPerUnit).Vec3(1.0))
	instance.MarkChanged()
	instance.Key = frame
	return
}

func (m *Manager) Bind() {
	m.cfg.Sheet.Bind()
}

func (m *Manager) Unbind() {
	m.cfg.Sheet.Unbind()
}

func (m *Manager) Delete() {
	m.cfg.Sheet.Delete()
	m.cfg.Sheet = nil
}

func (m *Manager) Render(camera *common.Camera) {
	m.Manager.Render(camera, m.cfg.Sheet)
}
