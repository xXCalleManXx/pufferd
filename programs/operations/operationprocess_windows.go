/*
 Copyright 2018 Padduck, LLC

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
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/command"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/download"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/mkdir"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/mojangdl"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/move"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/spongeforgedl"
	"github.com/pufferpanel/pufferd/programs/operations/ops/impl/writefile"
)

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

func loadOpModules() {
}
