package permissions

import (
	"encoding/json"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
)

var globalTracker PermissionTracker

func GetGlobal() PermissionTracker {
	if globalTracker == nil {
		data, err := ioutil.ReadFile(utils.JoinPath("data", "permissions.json"))
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
