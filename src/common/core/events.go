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

package core

import (
	"github.com/go-gl/glfw/v3.1/glfw"
)

type Events struct {
	window *glfw.Window
}

func newEvents(window *glfw.Window) (e *Events) {
	return &Events{
		window: window,
	}
}

func (e Events) Poll() {
	glfw.PollEvents()
}