package templates

import (
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
	"os"
)

const Folder = "data"

func CopyTemplates() {
	os.MkdirAll(Folder, os.ModeDir)

	data := Minecraft
	writeFile("minecraft", data)
}

func writeFile(name string, data string) {
	err := ioutil.WriteFile(utils.JoinPath(Folder, name+".json"), []byte(data), 0664)
	if err != nil {
		logging.Error("Error writing template "+name, err)
	}
}
