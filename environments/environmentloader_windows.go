package environments

import (
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
)

func LoadEnvironment(environmentType, folder, id string, environmentSection map[string]interface{}) Environment {
	switch environmentType {
	default:
		logging.Debugf("Loading server as standard")
		serverRoot := utils.JoinPath(folder, id)
		return &standard{RootDirectory: utils.GetStringOrDefault(environmentSection, "root", serverRoot), ConsoleBuffer: utils.CreateCache(), WSManager: utils.CreateWSManager()}
	}
}
