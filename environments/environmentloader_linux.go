package environments

import (
	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/cache"
	"github.com/pufferpanel/pufferd/utils"
)

func LoadEnvironment(environmentType, folder, id string, environmentSection map[string]interface{}) Environment {
	serverRoot := common.JoinPath(folder, id)
	rootDirectory := common.GetStringOrDefault(environmentSection, "root", serverRoot)
	cache := cache.CreateCache()
	wsManager := utils.CreateWSManager()
	switch environmentType {
	case "tty":
		logging.Debugf("Loading server as tty")
		t := &tty{standard: &standard{BaseEnvironment: &BaseEnvironment{}}}
		t.RootDirectory = rootDirectory
		t.ConsoleBuffer = cache
		t.WSManager = wsManager
		return t
	//case "docker":
	//	logging.Debugf("Loading server as tty")
	//	t := &docker{}
	//	t.RootDirectory = rootDirectory
	//	t.ConsoleBuffer = cache
	//	t.WSManager = wsManager
	//	return t
	default:
		logging.Debugf("Loading server as standard")
		s := &standard{BaseEnvironment: &BaseEnvironment{}}
		s.RootDirectory = rootDirectory
		s.ConsoleBuffer = cache
		s.WSManager = wsManager
		return s
	}
}
