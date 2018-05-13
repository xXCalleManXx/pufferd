package commands

import (
	"fmt"
	"github.com/pufferpanel/apufferi/config"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/uninstaller"
	"os"
	"strings"
)

func Uninstall(configPath string) {
	fmt.Println("This option will UNINSTALL pufferd, are you sure? Please enter \"yes\" to proceed [no]")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) == "yes" || strings.ToLower(response) == "y" {
		if os.Geteuid() != 0 {
			logging.Error("To uninstall pufferd you need to have sudo or root privileges")
		} else {
			config.Load(configPath)
			uninstaller.StartProcess()
			logging.Info("pufferd is now uninstalled.")
		}
	} else {
		logging.Info("Uninstall process aborted")
		logging.Info("Exiting")
	}
	return
}
