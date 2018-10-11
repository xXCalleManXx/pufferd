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
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/programs/operations/ops"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/command"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/download"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/mkdir"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/mojangdl"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/move"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/spongeforgedl"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/writefile"
)

var commandMapping map[string]ops.OperationFactory

func LoadOperations() {
	commandMapping = make(map[string]ops.OperationFactory)

	loadCoreModules()

	loadOpModules()
}

func GenerateProcess(directions []map[string]interface{}, environment environments.Environment, dataMapping map[string]interface{}, env map[string]string) OperationProcess {
	dataMap := make(map[string]interface{})
	for k, v := range dataMapping {
		dataMap[k] = v
	}

	//DEPRECATED: This will be removed in 1.4/2.0. This key should have been camelCased.
	dataMap["rootdir"] = environment.GetRootDirectory()

	dataMap["rootDir"] = environment.GetRootDirectory()
	operationList := make([]ops.Operation, 0)
	for _, mapping := range directions {

		factory := commandMapping[mapping["type"].(string)]

		opCreate := ops.CreateOperation{
			OperationArgs:        mapping,
			EnvironmentVariables: env,
			DataMap:              dataMapping,
		}

		op := factory.Create(opCreate)

		operationList = append(operationList, op)
	}
	return OperationProcess{processInstructions: operationList}
}

type OperationProcess struct {
	processInstructions []ops.Operation
}

func (p *OperationProcess) Run(env environments.Environment) (err error) {
	for p.HasNext() {
		err = p.RunNext(env)
		if err != nil {
			logging.Error("Error running process: ", err)
			break
		}
	}
	return
}

func (p *OperationProcess) RunNext(env environments.Environment) error {
	var op ops.Operation
	op, p.processInstructions = p.processInstructions[0], p.processInstructions[1:]
	err := op.Run(env)
	return err
}

func (p *OperationProcess) HasNext() bool {
	return len(p.processInstructions) != 0 && p.processInstructions[0] != nil
}

func loadCoreModules() {
	commandFactory := command.Factory
	commandMapping[commandFactory.Key()] = commandFactory

	downloadFactory := download.Factory
	commandMapping[downloadFactory.Key()] = downloadFactory

	mkdirFactory := mkdir.Factory
	commandMapping[mkdirFactory.Key()] = mkdirFactory

	moveFactory := move.Factory
	commandMapping[moveFactory.Key()] = moveFactory

	writeFileFactory := writefile.Factory
	commandMapping[writeFileFactory.Key()] = writeFileFactory

	mojangFactory := mojangdl.Factory
	commandMapping[mojangFactory.Key()] = mojangFactory

	spongeforgeDlFactory := spongeforgedl.Factory
	commandMapping[spongeforgeDlFactory.Key()] = spongeforgeDlFactory
}