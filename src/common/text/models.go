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
	"github.com/go-gl/mathgl/mgl32"
)

type ID int32

type Instance struct {
	renderIndex int
	tile        int
	position    mgl32.Vec3
	rotation    float32
	dirty       bool
	Text string
}

func (i *Instance) SetPosition(p mgl32.Vec3) {
	i.position = p
	i.dirty = true
}

func (i *Instance) SetRotation(r float32) {
	i.rotation = r
	i.dirty = true
}

type rendererInstance struct {
	model mgl32.Mat4
	tile  float32
}

type rendererData struct {
	Count     int
	Instances []rendererInstance
}
