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

package install

import (
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/programs/install/operations"
	"github.com/pufferpanel/pufferd/utils"
)

func GenerateInstallProcess(data *InstallSection, environment environments.Environment, dataMapping map[string]interface{}) InstallProcess {
	var directions = data.Commands
	datamap := make(map[string]interface{})
	for k, v := range dataMapping {
		datamap[k] = v.(map[string]interface{})["value"]
	}
	datamap["rootdir"] = environment.GetRootDirectory()
	ops := make([]operations.Operation, 0)
	for _, element := range directions {
		var mapping = element.(map[string]interface{})
		switch mapping["type"] {
		case "command":
			for _, element := range utils.ToStringArray(mapping["commands"]) {
				ops = append(ops, &operations.Command{Command: utils.ReplaceTokens(element, datamap), Environment: environment})
			}
		case "download":
			for _, element := range utils.ToStringArray(mapping["files"]) {
				ops = append(ops, &operations.Download{File: utils.ReplaceTokens(element, datamap), Environment: environment})
			}
		case "move":
			source := mapping["source"].(string)
			target := mapping["target"].(string)
			ops = append(ops, &operations.Move{SourceFile: source, TargetFile: target, Environment: environment})
		case "mkdir":
			target := mapping["target"].(string)
			ops = append(ops, &operations.Mkdir{TargetFile: target, Environment: environment})
		}
	}
	return InstallProcess{processInstructions: ops}
}

type InstallProcess struct {
	processInstructions []operations.Operation
}

func (p *InstallProcess) RunNext() error {
	var op operations.Operation
	op, p.processInstructions = p.processInstructions[0], p.processInstructions[1:]
	err := op.Run()
	return err
}

func (p *InstallProcess) HasNext() bool {
	return len(p.processInstructions) != 0 && p.processInstructions[0] != nil
}
