package operations

import (
	"path"
	"github.com/pufferpanel/apufferi/config"
	"io/ioutil"
	"os"
	"github.com/pufferpanel/apufferi/logging"
)

func loadOpsFromDir() {
	var directory = path.Join(config.GetStringOrDefault("", ""), "modules", "operations")

	files, err := ioutil.ReadDir(directory)
	if err != nil && os.IsNotExist(err) {
		return
	} else if err != nil {
		logging.Error("Error reading directory", err)
	}

	for _, file := range files {
		logging.Infof("Loading operation module: %s", file.Name())
	}
}
