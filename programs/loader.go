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

package programs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pufferpanel/pufferd/data/templates"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/install"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/permissions"
	"github.com/pufferpanel/pufferd/programs/types"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	programs     []Program = make([]Program, 0)
	ServerFolder string    = utils.JoinPath("data", "servers")
)

func LoadFromFolder() {
	os.Mkdir(ServerFolder, os.ModeDir)
	var programFiles, err = ioutil.ReadDir(ServerFolder)
	if err != nil {
		logging.Critical("Error reading from server data folder", err)
	}
	var data []byte
	var program Program
	for _, element := range programFiles {
		if element.IsDir() {
			continue
		}
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()))
		data, err = ioutil.ReadFile(utils.JoinPath(ServerFolder, element.Name()))
		if err != nil {
			logging.Error(fmt.Sprintf("Error loading server details from file (%s)", element.Name()), err)
			continue
		}
		program, err = LoadFromData(id, data)
		if err != nil {
			logging.Error(fmt.Sprintf("Error loading server details from json (%s)", element.Name()), err)
			continue
		}
		logging.Infof("Loaded server %s as %s", program.Id(), program.Name())
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
	var t = utils.GetStringOrDefault(pufferdData, "type", "")
	var installSection = getInstallSection(utils.GetMapOrNull(pufferdData, "install"))
	var runSection = utils.GetMapOrNull(pufferdData, "run")
	var environmentSection = utils.GetMapOrNull(runSection, "environment")
	var environment environments.Environment
	var defaultEnvType = "system"
	var environmentType = utils.GetStringOrDefault(environmentSection, "type", defaultEnvType)
	var permissions = permissions.Create(utils.GetMapOrNull(pufferdData, "permissions"))
	var dataSection = utils.GetMapOrNull(pufferdData, "data")
	dataCasted := make(map[string]string, len(dataSection))
	for key, value := range dataSection {
		dataCasted[key] = value.(string)
	}

	switch environmentType {
	case "system":
		serverRoot := utils.JoinPath(ServerFolder, id)
		environment = &environments.System{RootDirectory: utils.GetStringOrDefault(environmentSection, "root", serverRoot)}
	}

	switch t {
	case "java":
		var runBlock types.JavaRun
		if pufferdData["run"] == nil {
			runBlock = types.JavaRun{}
		} else {
			var stop = utils.GetStringOrDefault(runSection, "stop", "")
			var pre = utils.GetStringArrayOrNull(runSection, "pre")
			var post = utils.GetStringArrayOrNull(runSection, "post")
			var arguments = strings.Split(utils.GetStringOrDefault(runSection, "arguments", ""), " ")
			var enabled = utils.GetBooleanOrDefault(runSection, "enabled", true)
			var autostart = utils.GetBooleanOrDefault(runSection, "autostart", true)

			runBlock = types.NewJavaRun(stop, pre, post, arguments, dataCasted, enabled, autostart)
		}
		program = types.NewJavaProgram(id, runBlock, installSection, environment, permissions)
	}
	return
}

func Create(id string, serverType string, user string, data map[string]interface{}) {
	if GetFromCache(id) != nil {
		return
	}

	templateData, err := ioutil.ReadFile(utils.JoinPath(templates.Folder, serverType+".json"))

	var templateJson map[string]interface{}
	err = json.Unmarshal(templateData, &templateJson)
	segment := utils.GetMapOrNull(templateJson, "pufferd")

	if err != nil {
		logging.Error("Error reading template file for type "+serverType, err)
		return
	}

	if data != nil {
		segment["data"] = data
	}

	userPerms := make(map[string]interface{}, 0)
	userPerms[user] = append(make([]string, 0), ".*")
	segment["permissions"] = userPerms;

	templateData, _ = json.Marshal(templateJson)
	err = ioutil.WriteFile(utils.JoinPath(ServerFolder, id+".json"), templateData, 0644)

	if err != nil {
		logging.Error("Error writing server file", err)
		return
	}

	program, _ := LoadFromMapping(id, templateJson)
	programs = append(programs, program)
	program.Create()
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

	err = program.Destroy()
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

func getInstallSection(mapping map[string]interface{}) install.InstallSection {
	return install.Generate(
		utils.GetObjectArrayOrNull(mapping, "global"),
		utils.GetObjectArrayOrNull(mapping, "linux"),
		utils.GetObjectArrayOrNull(mapping, "mac"),
		utils.GetObjectArrayOrNull(mapping, "windows"))
}
