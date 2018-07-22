package operations

import (
	"github.com/pufferpanel/apufferi/config"
	"github.com/pufferpanel/apufferi/logging"
	"io/ioutil"
	"os"
	"path"
	"plugin"
	"github.com/pufferpanel/pufferd/programs/operations/ops"
	"reflect"
)

func loadOpModules() {
	var directory = path.Join(config.GetStringOrDefault("dataFolder", ""), "modules", "operations")

	files, err := ioutil.ReadDir(directory)
	if err != nil && os.IsNotExist(err) {
		return
	} else if err != nil {
		logging.Error("Error reading directory", err)
	}

	for _, file := range files {
		logging.Infof("Loading operation module: %s", file.Name())
		p, e := plugin.Open(path.Join(directory, file.Name()))
		if err != nil {
			logging.Error("Unable to load module", e)
			continue
		}

		factory, e := p.Lookup("Factory")
		if err != nil {
			logging.Error("Unable to load module", e)
			continue
		}

		fty, ok := factory.(ops.OperationFactory)
		if !ok {
			logging.Errorf("Expected OperationFactory, but found %s", reflect.TypeOf(factory).Name())
			continue
		}

		commandMapping[fty.Key()] = fty

		logging.Infof("Loaded operation module: %s", fty.Key())
	}
}
