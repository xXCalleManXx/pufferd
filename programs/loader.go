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
	"github.com/pufferpanel/pufferd/programs/types"
	"io/ioutil"
	"strings"
	"path/filepath"
	"github.com/pufferpanel/pufferd/logging"
	"fmt"
	"os"
)

const (
	serverFolder = "servers";
)

var (
	programs []Program = make([]Program, 5);
)

func LoadFromFolder() {
	os.Mkdir(serverFolder, os.ModeDir);
	var programFiles, err = ioutil.ReadDir(serverFolder);
	if (err != nil) {
		logging.Critical("Error reading from server data folder", err);
	}
	var data []byte;
	var program Program;
	for _, element := range programFiles {
		id := strings.TrimSuffix(element.Name(), filepath.Ext(element.Name()));
		data, err = ioutil.ReadFile(joinPath(serverFolder, element.Name()));
		if (err != nil) {
			logging.Error(fmt.Sprintf("Error loading server details (%s)", element.Name()), err);
			continue;
		}
		program, err = LoadProgramFromData(id, data);
		if (err != nil) {
			logging.Error(fmt.Sprintf("Error loading server details (%s)", element.Name()), err);
			continue;
		}
		logging.Infof("Loaded server %s as %s", program.Id(), program.Name());
		programs = append(programs, program);
	}
}

func GetProgram(id string) (program Program, err error) {
	program = getFromCache(id);
	if (program == nil) {
		program, err = LoadProgram(id);
	}
	return;
}

func LoadProgram(id string) (program Program, err error) {
	var data []byte;
	data, err = ioutil.ReadFile(joinPath(serverFolder, id + ".json"));
	program, err = LoadProgramFromData(id, data);
	return;
}

func LoadProgramFromData(id string, source []byte) (program Program, err error) {
	var data map[string]interface{};
	err = json.Unmarshal(source, &data);
	if (err != nil) {
		return;
	}
	var pufferdData = GetMapOrNull(data, "pufferd");
	var t = GetStringOrDefault(pufferdData, "type", nil);
	var install = GetInstallSection(GetMapOrNull(pufferdData, "install"));
	switch(t) {
	case "java":
		var runBlock types.JavaRun;
		if (pufferdData["run"] == nil) {
			runBlock = types.JavaRun{};
		} else {
			var runSection = GetMapOrNull(pufferdData, "run");
			var stop = GetStringOrDefault(runSection, "stop", nil);
			var pre = GetStringArrayOrNull(runSection, "pre");
			var post = GetStringArrayOrNull(runSection, "post");
			var arguments = GetStringOrDefault(runSection, "arguments", nil);
			var enabled = GetBooleanOrDefault(runSection, "enabled", true);

			runBlock = types.JavaRun{Stop: stop, Pre: pre, Post: post, Arguments: arguments, Enabled: enabled};
		}
		program = types.NewJavaProgram(id, runBlock, install);
	}
	return;
}

func GetInstallSection(data map[string]interface{}) types.JavaInstall {
	var install = types.JavaInstall{};
	install.Files = GetStringArrayOrNull(data, "files");
	install.Pre = GetStringArrayOrNull(data, "pre");
	install.Post = GetStringArrayOrNull(data, "post");
	return install;
}

func GetStringOrDefault(data map[string]interface{}, key string, def *string) string {
	if (data == nil) {
		return *def;
	}
	var section = data[key];
	if (section == nil) {
		return *def;
	} else {
		return section.(string);
	}
}

func GetBooleanOrDefault(data map[string]interface{}, key string, def bool) bool {
	if (data == nil) {
		return def;
	}
	var section = data[key];
	if (section == nil) {
		return def;
	} else {
		return section.(bool);
	}
}

func GetMapOrNull(data map[string]interface{}, key string) map[string]interface{} {
	if (data == nil) {
		return (map[string]interface{})(nil);
	}
	var section = data[key];
	if (section == nil) {
		return (map[string]interface{})(nil);
	} else {
		return section.(map[string]interface{});
	}
}

func GetStringArrayOrNull(data map[string]interface{}, key string) []string {
	if (data == nil) {
		return ([]string)(nil);
	}
	var section = data[key];
	if (section == nil) {
		return ([]string)(nil);
	} else {
		var sec = section.([]interface{});
		var newArr = make([]string, len(sec));
		for i := 0; i < len(sec); i++ {
			newArr[i] = sec[i].(string);
		}
		return newArr;
	}
}

func joinPath(paths ...string) string {
	return strings.Join(paths, string(filepath.Separator));
}

func getFromCache(id string) (Program) {
	for _, element := range programs {
		if (element.Id() == id) {
			return element;
		}
	}
	return nil;
}