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

package resources

import (
	"github.com/kurrik/opengl-benchmarks/common"
	"github.com/kurrik/opengl-benchmarks/common/render"
	"github.com/kurrik/opengl-benchmarks/common/sprites"
)

type Loader interface {
	Load(m *Manager) (err error)
}

type Manager struct {
	geometry map[string]*render.Geometry
	textures map[string]*common.Texture
	sheets   map[string]*sprites.Sheet
}

func NewManager() *Manager {
	return &Manager{
		geometry: map[string]*render.Geometry{},
		textures: map[string]*common.Texture{},
		sheets:   map[string]*sprites.Sheet{},
	}
}

func (m *Manager) SetGeometry(key string, geometry *render.Geometry) {
	var (
		exists bool
		g      *render.Geometry
	)
	if g, exists = m.geometry[key]; exists {
		g.Delete()
	}
	m.geometry[key] = geometry
}

func (m *Manager) GetGeometry(key string) *render.Geometry {
	return m.geometry[key]
}

func (m *Manager) SetTexture(key string, texture *common.Texture) {
	var (
		exists bool
		t      *common.Texture
	)
	if t, exists = m.textures[key]; exists {
		t.Delete()
	}
	m.textures[key] = texture
}

func (m *Manager) SetSheet(key string, sheet *sprites.Sheet) {
	var (
		exists bool
		s      *sprites.Sheet
	)
	if s, exists = m.sheets[key]; exists {
		s.Delete()
	}
	m.sheets[key] = sheet
}

func (m *Manager) GetSheet(key string) *sprites.Sheet {
	return m.sheets[key]
}

func (m *Manager) Delete() {
	for _, g := range m.geometry {
		g.Delete()
	}
	for _, s := range m.sheets {
		s.Delete()
	}
	for _, t := range m.textures {
		t.Delete()
	}
}

func (m *Manager) Load(loaders ...Loader) (err error) {
	var (
		loader Loader
	)
	for _, loader = range loaders {
		if err = loader.Load(m); err != nil {
			return
		}
	}
	return
}
