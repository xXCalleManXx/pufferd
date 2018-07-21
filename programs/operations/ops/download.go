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

package ops

import (
	"github.com/cavaliercoder/grab"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/apufferi/common"
)

type Download struct {
	Files        []string
}

func (d Download) Run(env environments.Environment) error {
	for _, file := range d.Files {
		logging.Debugf("Download file from %s to %s", file, env.GetRootDirectory())
		env.DisplayToConsole("Downloading file %s\n", file)
		_, err := grab.Get(env.GetRootDirectory(), file)
		if err != nil {
			return err
		}
	}
	return nil
}

type DownloadOperationFactory struct {
}

func (of DownloadOperationFactory) 	Create(op CreateOperation) Operation {
	files := common.ToStringArray(op.OperationArgs["files"])
	return &Download{Files: files}
}

func (of DownloadOperationFactory) Key() string {
	return "download"
}