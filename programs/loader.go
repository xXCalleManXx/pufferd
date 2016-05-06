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
)

func LoadServer(id string) (program Program, err error) {
	var data []byte;
	return LoadServerFromData(data);
}

func LoadServerFromData(source []byte) (program Program, err error) {
	var data map[string]interface{};
	err = json.Unmarshal(source, &data);
	if (err != nil) {
		return;
	}
	var pufferdData = GetMapOrNull(data, "pufferd");
	var t = GetStringOrNull(pufferdData, "type");
	switch(t) {
	case "java":
		var runBlock types.JavaRun;
		if (pufferdData["run"] == nil) {
			runBlock = types.JavaRun{};
		} else {
			var runSection = GetMapOrNull(pufferdData, "run");
			var stop = GetStringOrNull(runSection, "stop");
			var pre = GetStringArrayOrNull(runSection, "pre");
			var post = GetStringArrayOrNull(runSection, "post");
			var arguments = GetStringOrNull(runSection, "arguments");

			runBlock = types.JavaRun{Stop: stop, Pre: pre, Post: post, Arguments: arguments};
		}
		program = &types.Java{Run: runBlock};
	}
	return;
}

func GetStringOrNull(data map[string]interface{}, key string) string {
	if (data == nil) {
		return "";
	}
	var section = data[key];
	if (section == nil) {
		return "";
	} else {
		return section.(string);
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