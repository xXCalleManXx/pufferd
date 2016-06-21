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

package templates

import (
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
	"os"
)

var Folder = utils.JoinPath("data", "templates")

func CopyTemplates() {
	os.MkdirAll(Folder, os.ModeDir)

	data := Spigot
	writeFile("spigot", data)

	data = CraftbukkitBySpigot
	writeFile("fakecraftbukkit", data)

	data = Vanilla
	writeFile("vanilla", data)
}

func writeFile(name string, data string) {
	err := ioutil.WriteFile(utils.JoinPath(Folder, name+".json"), []byte(data), 0664)
	if err != nil {
		logging.Error("Error writing template "+name, err)
	}
}
