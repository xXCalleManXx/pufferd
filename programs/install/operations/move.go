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
	"path/filepath"

	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
)

type Move struct {
	SourceFile  string
	TargetFile  string
	Environment environments.Environment
}

func (m *Move) Run() error {
	source := utils.JoinPath(m.Environment.GetRootDirectory(), m.SourceFile)
	target := utils.JoinPath(m.Environment.GetRootDirectory(), m.TargetFile)
	result, valid := validateMove(source, target)
	if !valid {
		return nil
	}
	for k, v := range result {
		logging.Debugf("Moving file from %s to %s", source, target)
		m.Environment.DisplayToConsole("Moving file from %s to %s", m.SourceFile, m.TargetFile)
		err := os.Rename(k, v)
		if err != nil {
			logging.Error("Error moving file", err)
		}
	}
	return nil
}

func validateMove(source string, target string) (result map[string]string, valid bool) {
	result = make(map[string]string)
	sourceFiles, _ := filepath.Glob(source)
	info, err := os.Stat(target)

	if err != nil {
		if os.IsNotExist(err) && len(sourceFiles) > 1 {
			logging.Error("Target folder does not exist", err)
			valid = false
			return
		} else if !os.IsNotExist(err) {
			valid = false
			logging.Error("Error reading target file on move", err)
			return
		}
	} else if info.IsDir() && len(sourceFiles) > 1 {
		logging.Error("Cannot move multiple files to single file target")
		valid = false
		return
	}

	if info != nil && info.IsDir() {
		for _, v := range sourceFiles {
			_, fileName := filepath.Split(v)
			result[v] = filepath.Join(target, fileName)
		}
	} else {
		for _, v := range sourceFiles {
			result[v] = target
		}
	}
	valid = true
	return
}
