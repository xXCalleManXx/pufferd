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

package install

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/data"
)

func Install(configPath string, authRoot string, authToken string) {
	if authRoot == "" {
		logging.Error("Authorization server root not passed")
		os.Exit(1)
	}

	if authToken == "" {
		logging.Error("Authorization token not passed")
		os.Exit(1)
	}

	config := data.CONFIG

	replacements := make(map[string]interface{})
	replacements["authUrl"] = strings.TrimSuffix(authRoot, "/")
	replacements["authToken"] = authToken

	configData := []byte(common.ReplaceTokens(config, replacements))

	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, configData, "", "  ")
	if err != nil {
		logging.Error("Error writing new config")
		os.Exit(1)
	}

	err = ioutil.WriteFile(configPath, prettyJson.Bytes(), 0664)

	if err != nil {
		logging.Error("Error writing new config")
		os.Exit(1)
	}

	logging.Info("Config saved")
}
