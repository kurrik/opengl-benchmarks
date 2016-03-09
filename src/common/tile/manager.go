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

package tile

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
)

type Config struct {
	MaxInstances uint32
}

type Manager struct {
	cfg               Config
	renderer          *Renderer
	Instances         InstanceList
	rendererInstances []rInstance
	count             int
}

func NewManager(cfg Config) (mgr *Manager, err error) {
	mgr = &Manager{
		cfg:               cfg,
		rendererInstances: make([]rInstance, cfg.MaxInstances),
		count:             0,
	}
	if mgr.renderer, err = NewRenderer(); err != nil {
		return
	}
	return
}

func (m *Manager) CreateInstance() (inst *Instance, err error) {
	inst = NewInstance()
	inst.SetPosition(mgl32.Vec3{0, 0, 0})
	inst.SetScale(mgl32.Vec3{1.0, 1.0, 1.0})
	inst.SetRotation(0)
	inst.Tile = 0
	m.Instances.Prepend(inst)
	m.count += 1
	return
}

func (m *Manager) Bind() {
	m.renderer.Bind()
}

func (m *Manager) Unbind() {
	m.renderer.Unbind()
}

func (m *Manager) Delete() {
	m.renderer.Delete()
}

func (m *Manager) Render(camera *common.Camera, sheet *Sheet) {
	var (
		inst  *Instance
		rinst *rInstance
		scale mgl32.Vec3
		index int
	)
	scale = mgl32.Vec3{
		1.0 / camera.PxPerUnit.X(),
		1.0 / camera.PxPerUnit.Y(),
		1.0,
	}
	index = 0
	inst = m.Instances.Head()
	for inst != nil {
		if uint32(index) >= m.cfg.MaxInstances {
			m.renderer.Render(camera, index, m.rendererInstances, sheet)
			index = 0
		}
		inst.SetScale(scale)
		rinst = &m.rendererInstances[index]
		rinst.tile = float32(inst.Tile)
		rinst.model = inst.GetModel()
		index++
		inst = inst.Next()
	}
	m.renderer.Render(camera, index, m.rendererInstances, sheet)
}
