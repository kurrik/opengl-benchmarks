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

package sprites

import (
	"github.com/kurrik/opengl-benchmarks/common/render"
)

type Manager struct {
	pixelsPerUnit float32
}

func NewManager(pixelsPerUnit float32) *Manager {
	return &Manager{
		pixelsPerUnit: pixelsPerUnit,
	}
}

func (m *Manager) SetFrame(instance *render.Instance, sheet *Sheet, frame string) (err error) {
	var s *Sprite
	if instance == nil {
		return // No error
	}
	if s, err = sheet.Sprite(frame); err != nil {
		return
	}
	instance.Frame = s.Index()
	instance.SetScale(s.WorldDimensions(m.pixelsPerUnit).Vec3(1.0))
	instance.MarkChanged()
	instance.Key = frame
	return
}
