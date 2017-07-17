// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loggers

import (
	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/core/logging/structure"
)

// Treat duplicate key-values as consecutive entries in a vector-valued lookup
type vectorValuedLogger struct {
	logger log.Logger
}

func (vvl *vectorValuedLogger) Log(keyvals ...interface{}) error {
	return vvl.logger.Log(structure.Vectorise(keyvals)...)
}

func VectorValuedLogger(logger log.Logger) log.Logger {
	return &vectorValuedLogger{logger: logger}
}