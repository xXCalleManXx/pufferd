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
	"bytes"
	"encoding/json"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
	"os"
	"github.com/pufferpanel/pufferd/config"
)

var Folder = config.GetOrDefault("templatefolder", utils.JoinPath("data", "templates"))

func CopyTemplates() {
	os.MkdirAll(Folder, 0755)

	data := Spigot
	writeFile("spigot", data)
	writeFile("minecraft", data)

	data = Bungeecord
	writeFile("bungeecord", data)

	data = CraftbukkitBySpigot
	writeFile("fakecraftbukkit", data)

	data = Vanilla
	writeFile("vanilla", data)

	data = Forge
	writeFile("forge", data)

	data = Sponge
	writeFile("spongeforge", data)

	data = SRCDS
	writeFile("srcds", data)

	data = TF2
	writeFile("tf2", data)
}

func writeFile(name string, data string) {
	jsonData := []byte(data)
	var prettyJson bytes.Buffer
	json.Indent(&prettyJson, jsonData, "", "  ")
	err := ioutil.WriteFile(utils.JoinPath(Folder, name+".json"), prettyJson.Bytes(), 0664)
	if err != nil {
		logging.Error("Error writing template "+name, err)
	}
}
