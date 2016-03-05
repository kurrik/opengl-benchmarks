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
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kurrik/opengl-benchmarks/common"
)

type ManagerConfig struct {
	MaxInstances uint32
}

type Manager struct {
	cfg               ManagerConfig
	nextID            InstanceID
	renderer          *Renderer
	Instances         map[InstanceID]*TileInstance
	rendererInstances []rInstance
	count             int
}

func NewManager(cfg ManagerConfig) (mgr *Manager, err error) {
	mgr = &Manager{
		cfg:               cfg,
		Instances:         map[InstanceID]*TileInstance{},
		rendererInstances: make([]rInstance, cfg.MaxInstances),
		count:             0,
	}
	if mgr.renderer, err = NewRenderer(); err != nil {
		return
	}
	return
}

func (m *Manager) CreateInstance() (id InstanceID, err error) {
	if uint32(m.count) >= m.cfg.MaxInstances {
		err = fmt.Errorf("Max instances reached")
		return
	}
	id = m.nextID
	m.Instances[id] = &TileInstance{
		renderIndex: m.count,
		position:    mgl32.Vec3{0, 0, 0},
		rotation:    0,
		Tile:        0,
		Dirty:       true,
	}
	m.nextID += 1
	m.count += 1
	return
}

func (m *Manager) GetInstance(id InstanceID) (inst *TileInstance, err error) {
	var (
		exists bool
	)
	if inst, exists = m.Instances[id]; !exists {
		err = fmt.Errorf("Invalid text instance ID: %v", id)
		return
	}
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
		inst  *TileInstance
		rinst *rInstance
		scale mgl32.Mat4
		rot   mgl32.Mat4
		trans mgl32.Mat4
	)
	scale = mgl32.Scale3D(
		1.0/camera.PxPerUnit.X(),
		1.0/camera.PxPerUnit.Y(),
		1.0,
	)
	for _, inst = range m.Instances {
		if inst.Dirty {
			rinst = &m.rendererInstances[inst.renderIndex]
			rinst.tile = float32(inst.Tile)
			rot = mgl32.HomogRotate3DZ(mgl32.DegToRad(inst.rotation))
			trans = mgl32.Translate3D(
				inst.position.X(),
				inst.position.Y(),
				inst.position.Z(),
			)
			rinst.model = trans.Mul4(rot).Mul4(scale)
			inst.Dirty = false
		}
	}
	m.renderer.Render(camera, m.count, m.rendererInstances, sheet)
}
