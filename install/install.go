package install

import (
	"github.com/pufferpanel/pufferd/logging"
	"os"
	"github.com/pufferpanel/pufferd/data"
	"strings"
	"github.com/pufferpanel/pufferd/utils"
	"bytes"
	"encoding/json"
	"io/ioutil"
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
	replacements["authurl"] = strings.TrimSuffix(authRoot, "/")
	replacements["authtoken"] = authToken

	configData := []byte(utils.ReplaceTokens(config, replacements))

	var prettyJson bytes.Buffer
	json.Indent(&prettyJson, configData, "", "  ")
	err := ioutil.WriteFile(configPath, prettyJson.Bytes(), 0664)

	if err != nil {
		logging.Error("Error writing new config")
		os.Exit(1)
	}

	logging.Info("Config saved")

	logging.Info("Attempting to install service")
	InstallService()
}
