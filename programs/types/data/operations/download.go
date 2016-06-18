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
	"fmt"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/logging"
	"os"
	"path/filepath"
)

type Download struct {
	File        string
	Environment environments.Environment
}

func (d *Download) Run() {
	fmt.Println("Downloading file: " + d.File)
	_, fileName := filepath.Split(d.File)
	logging.Debug(os.Getenv("path"))
	logging.Debug(os.Getenv("PATH"))
	d.Environment.Execute("curl", []string{"-o", fileName, d.File})
}
