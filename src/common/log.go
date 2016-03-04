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

package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Logger struct {
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

func NewLogger(level string) (out *Logger, err error) {
	var (
		flag   int = log.Lshortfile | log.LstdFlags | log.Lmicroseconds
		nilLog     = log.New(ioutil.Discard, "", flag)
		writer     = os.Stderr
	)
	out = &Logger{
		Debug: nilLog,
		Info:  nilLog,
		Error: nilLog,
	}
	switch strings.ToLower(level) {
	case "debug":
		out.Debug = log.New(writer, "D ", flag)
		fallthrough
	case "info":
		out.Info = log.New(writer, "I ", flag)
		fallthrough
	case "error":
		out.Error = log.New(writer, "E ", flag)
	default:
		err = fmt.Errorf("Invalid log level: %v", level)
	}
	return
}
