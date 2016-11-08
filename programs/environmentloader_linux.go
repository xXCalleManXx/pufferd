package programs

import (
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
)

func LoadEnvironment(environmentType, folder, id string, environmentSection map[string]interface{}) environments.Environment {
	switch environmentType {
	case "tty":
		logging.Debugf("Loading server as tty")
		serverRoot := utils.JoinPath(folder, id)
		return &environments.Tty{RootDirectory: utils.GetStringOrDefault(environmentSection, "root", serverRoot), ConsoleBuffer: utils.CreateCache(), WSManager: utils.CreateWSManager()}
	default:
		logging.Debugf("Loading server as standard")
		serverRoot := utils.JoinPath(folder, id)
		return &environments.Standard{RootDirectory: utils.GetStringOrDefault(environmentSection, "root", serverRoot), ConsoleBuffer: utils.CreateCache(), WSManager: utils.CreateWSManager()}
	}
}
