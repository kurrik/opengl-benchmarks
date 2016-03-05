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
	"github.com/kurrik/opengl-benchmarks/common"
)

type Sheet struct {
	uniformData *Uniform
	texture     *common.Texture
	keys        map[string]int
}

func NewSheet() *Sheet {
	return &Sheet{
		uniformData: NewUniform(),
		keys:        map[string]int{},
	}
}

func (s *Sheet) AddTile(key string, tile Tile) {
	s.keys[key] = s.uniformData.AppendTile(tile)
}
