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
