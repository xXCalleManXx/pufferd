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

package permissions

import (
	"encoding/json"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
	"os"
)

var globalTracker PermissionTracker

func GetGlobal() PermissionTracker {
	if globalTracker == nil {
		file := utils.JoinPath("data", "permissions.json")
		if _, err := os.Stat(file); os.IsNotExist(err) {
			ioutil.WriteFile(file, []byte("{}"), 0664)
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			logging.Error("Error loading global permissions", err)
		}
		var mapped map[string]interface{}
		err = json.Unmarshal(data, &mapped)
		if err != nil {
			logging.Error("Error loading global permissions", err)
		}
		globalTracker = Create(mapped)
	}
	return globalTracker
}
