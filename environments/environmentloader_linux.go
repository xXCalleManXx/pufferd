package environments

import (
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
)

func LoadEnvironment(environmentType, folder, id string, environmentSection map[string]interface{}) Environment {
	serverRoot := utils.JoinPath(folder, id)
	rootDirectory := utils.GetStringOrDefault(environmentSection, "root", serverRoot)
	switch environmentType {
	case "tty":
		logging.Debugf("Loading server as tty")
		return &Tty{RootDirectory: rootDirectory, ConsoleBuffer: utils.CreateCache(), WSManager: utils.CreateWSManager()}
	//case "docker":
	//	logging.Debugf("Loading server as docker")
	//	netBindings := make([]string, 0)
	//	return &Docker{ContainerId: id, RootDirectory: rootDirectory, ConsoleBuffer: utils.CreateCache(), WSManager: utils.CreateWSManager(), NetworkBindings: netBindings}
	default:
		logging.Debugf("Loading server as standard")
		return &Standard{RootDirectory: rootDirectory, ConsoleBuffer: utils.CreateCache(), WSManager: utils.CreateWSManager()}
	}
}
