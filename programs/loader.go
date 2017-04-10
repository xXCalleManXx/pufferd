/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 distributed under the License is distributed on an "AS IS" BASIS,
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package programs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pufferpanel/pufferd/config"
	"github.com/pufferpanel/pufferd/data/templates"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs/install"
	"github.com/pufferpanel/pufferd/utils"
)

var (
	programs     []Program = make([]Program, 0)
	ServerFolder string
)

func Initialize() {
	ServerFolder = config.GetOrDefault("serverfolder", utils.JoinPath("data", "servers"))
}

func LoadFromFolder() {
	os.Mkdir(ServerFolder, 0755)
	var programFiles, err = ioutil.ReadDir(ServerFolder)
	if err != nil {
		logging.Critical("Error reading from server data folder", err)
	}
	var program Program
	for _, element := range programFiles {
		if element.IsDir() {
			continue
		}
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		program, err = Load(id)
		if err != nil {
			logging.Error(fmt.Sprintf("Error loading server details from json (%s)", element.Name()), err)
			continue
		}
		logging.Infof("Loaded server %s", program.Id())
		programs = append(programs, program)
	}
}

func Get(id string) (program Program, err error) {
	program = GetFromCache(id)
	if program == nil {
		program, err = Load(id)
	}
	return
}

func GetAll() []Program {
	return programs
}

func Load(id string) (program Program, err error) {
	var data []byte
	data, err = ioutil.ReadFile(utils.JoinPath(ServerFolder, id+".json"))
	if len(data) == 0 || err != nil {
		return
	}

	program, err = LoadFromData(id, data)
	return
}

func LoadFromData(id string, source []byte) (program Program, err error) {
	var data map[string]interface{}
	err = json.Unmarshal(source, &data)
	if err != nil {
		return
	}
	program, err = LoadFromMapping(id, data)
	return
}

func LoadFromMapping(id string, source map[string]interface{}) (program Program, err error) {
	var pufferdData = utils.GetMapOrNull(source, "pufferd")
	var installSection = getInstallSection(utils.GetMapOrNull(pufferdData, "install"))
	var runSection = utils.GetMapOrNull(pufferdData, "run")
	var environmentSection = utils.GetMapOrNull(pufferdData, "environment")
	var environment environments.Environment
	var dataSection = utils.GetMapOrNull(pufferdData, "data")
	dataCasted := make(map[string]interface{}, len(dataSection))
	for key, value := range dataSection {
		dataCasted[key] = value
	}

	var environmentType string
	if environmentSection == nil {
		environmentType = "standard"
	} else {
		environmentType = utils.GetStringOrDefault(environmentSection, "type", "standard")
	}

	logging.Debugf("Loading server as %s", environmentType)

	environment = environments.LoadEnvironment(environmentType, ServerFolder, id, environmentSection)

	var runBlock Runtime
	if pufferdData["run"] == nil {
		runBlock = Runtime{}
	} else {
		var stop = utils.GetStringOrDefault(runSection, "stop", "")
		var pre = utils.GetStringArrayOrNull(runSection, "pre")
		var post = utils.GetStringArrayOrNull(runSection, "post")
		var arguments = utils.GetStringArrayOrNull(runSection, "arguments")
		var enabled = utils.GetBooleanOrDefault(runSection, "enabled", true)
		var autostart = utils.GetBooleanOrDefault(runSection, "autostart", true)
		var program = utils.GetStringOrDefault(runSection, "program", "")
		runBlock = Runtime{Stop: stop, Pre: pre, Post: post, Arguments: arguments, Enabled: enabled, AutoStart: autostart, Program: program}
	}
	program = &programData{Data: dataCasted, Identifier: id, RunData: runBlock, InstallData: installSection, Environment: environment}
	return
}

func Create(id string, serverType string, data map[string]interface{}) bool {
	if GetFromCache(id) != nil {
		return false
	}

	templateData, err := ioutil.ReadFile(utils.JoinPath(templates.Folder, serverType+".json"))
	if err != nil {
		logging.Error("Error reading template file for type "+serverType, err)
		return false
	}

	var templateJson map[string]interface{}
	err = json.Unmarshal(templateData, &templateJson)
	segment := utils.GetMapOrNull(templateJson, "pufferd")

	if err != nil {
		logging.Error("Error reading template file for type "+serverType, err)
		return false
	}

	if data != nil {
		var mapper map[string]interface{}
		mapper = segment["data"].(map[string]interface{})
		for k, v := range data {
			if mapper[k] == nil {
				newMap := make(map[string]interface{})
				newMap["value"] = v
				newMap["desc"] = "No description"
				newMap["display"] = k
				newMap["required"] = false
				newMap["internal"] = true
				mapper[k] = newMap
			} else {
				mapper[k].(map[string]interface{})["value"] = v
			}
		}
		segment["data"] = mapper
	}

	f, err := os.Create(utils.JoinPath(ServerFolder, id+".json"))

	if err != nil {
		logging.Error("Error writing server file", err)
		return false
	}

	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(templateJson)

	if err != nil {
		logging.Error("Error writing server file", err)
		return false
	}

	program, _ := LoadFromMapping(id, templateJson)
	programs = append(programs, program)
	program.Create()
	return true
}

func Delete(id string) (err error) {
	var index int
	var program Program
	for i, element := range programs {
		if element.Id() == id {
			program = element
			index = i
			break
		}
	}
	if program == nil {
		return
	}

	if program.IsRunning() {
		err = program.Stop()
		if err != nil {
			return err
		}
	}

	err = program.Destroy()
	if err != nil {
		return err
	}
	os.Remove(utils.JoinPath(ServerFolder, program.Id()+".json"))
	programs = append(programs[:index], programs[index+1:]...)
	return
}

func GetFromCache(id string) Program {
	for _, element := range programs {
		if element.Id() == id {
			return element
		}
	}
	return nil
}

func Save(id string) (err error) {
	program := GetFromCache(id)
	if program == nil {
		err = errors.New("No server with given id")
		return
	}
	err = program.Save(utils.JoinPath(ServerFolder, id+".json"))
	return
}

func Reload(id string) error {
	oldPg, err := Get(id)
	if err != nil {
		return err
	}

	var newPg Program

	newPg, err = Load(id)
	if err != nil {
		return err
	}

	oldPg.Reload(newPg)
	return nil
}

func GetPlugins() map[string]interface{} {

	temps, _ := ioutil.ReadDir(templates.Folder)

	mapping := make(map[string]interface{})

	for _, element := range temps {
		if element.IsDir() {
			continue
		}
		name := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		templateData, _ := ioutil.ReadFile(utils.JoinPath(templates.Folder, name+".json"))

		var templateJson map[string]interface{}
		err := json.Unmarshal(templateData, &templateJson)
		if err != nil {
			logging.Error("Malformed json for program "+element.Name(), err)
			continue
		}
		segment := utils.GetMapOrNull(templateJson, "pufferd")
		dataSec := make(map[string]interface{})
		dataSec["variables"] = segment["data"].(map[string]interface{})
		dataSec["display"] = segment["display"]
		mapping[name] = dataSec
	}

	return mapping
}

func getInstallSection(mapping map[string]interface{}) install.InstallSection {
	return install.InstallSection{
		Commands: utils.GetObjectArrayOrNull(mapping, "commands"),
	}
}
