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

type TextID int32

type managerInstance struct {
	packedIndex int
	position mgl32.Vec2
	rotation float32
}

type rendererInstance struct {
	model mgl32.Mat4
	tile  float32
}

type rendererData struct {
	Instances []rendererInstance
}
