package environments

import (
	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/cache"
	"github.com/pufferpanel/pufferd/utils"
)

func LoadEnvironment(environmentType, folder, id string, environmentSection map[string]interface{}) Environment {
	switch environmentType {
	default:
		logging.Debugf("Loading server as standard")
		serverRoot := common.JoinPath(folder, id)
		return &standard{RootDirectory: common.GetStringOrDefault(environmentSection, "root", serverRoot), ConsoleBuffer: cache.CreateCache(), WSManager: utils.CreateWSManager()}
	}
}
