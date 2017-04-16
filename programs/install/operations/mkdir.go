/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package operations

import (
	"os"

	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/pufferpanel/pufferd/logging"
)

type Mkdir struct {
	TargetFile  string
	Environment environments.Environment
}

func (m *Mkdir) Run() error {
	logging.Debugf("Making directory: %s\n", m.TargetFile)
	m.Environment.DisplayToConsole("Creating directory: %s\n", m.TargetFile)
	target := utils.JoinPath(m.Environment.GetRootDirectory(), m.TargetFile)
	return os.MkdirAll(target, 0755)
}
