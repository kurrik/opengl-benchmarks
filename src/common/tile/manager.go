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
	"github.com/kurrik/opengl-benchmarks/common/render"
	"github.com/kurrik/opengl-benchmarks/common/sprites"
)

type Manager struct {
	Instances render.InstanceList
	count     int
	geometry  *render.Geometry
	renderer  *render.Renderer
}

func NewManager(renderer *render.Renderer) (mgr *Manager, err error) {
	mgr = &Manager{
		count:    0,
		renderer: renderer,
		geometry: render.NewGeometryFromPoints(render.Square),
	}
	return
}

func (m *Manager) CreateInstance() (inst *render.Instance, err error) {
	inst = render.NewInstance()
	inst.SetPosition(mgl32.Vec3{0, 0, 0})
	inst.SetScale(mgl32.Vec3{1.0, 1.0, 1.0})
	inst.SetRotation(0)
	inst.Frame = 0
	m.Instances.Prepend(inst)
	m.count += 1
	return
}

func (m *Manager) Render(camera *common.Camera, sheet sprites.UniformBufferSheet) {
	m.renderer.Render(camera, sheet, m.geometry, &m.Instances)
}
