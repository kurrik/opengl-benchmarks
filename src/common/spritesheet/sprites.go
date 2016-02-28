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
	"io/ioutil"
	"path"
)

type Sprites struct {
	Sheet   *Spritesheet
	Texture *common.Texture
}

func NewSprites(jsonPath string, pxPerUnit float32) (sprites *Sprites, err error) {
	var (
		data        []byte
		dir         = path.Dir(jsonPath)
		texture     *common.Texture
		spritesheet *Spritesheet
	)
	if data, err = ioutil.ReadFile(jsonPath); err != nil {
		return
	}
	if spritesheet, err = ParseTexturePackerJSONArrayString(
		string(data),
		pxPerUnit,
	); err != nil {
		return
	}
	if texture, err = common.LoadTexture(
		path.Join(dir, spritesheet.GetTexturePath()),
		common.SmoothingNearest,
	); err != nil {
		return
	}
	sprites = &Sprites{
		Sheet:   spritesheet,
		Texture: texture,
	}
	return
}
