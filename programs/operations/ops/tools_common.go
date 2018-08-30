package ops

import (
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/environments"
	"io"
	"net/http"
	"os"
	"path"
)

func downloadFile(url, fileName string, env environments.Environment) error {
	target, err := os.Create(path.Join(env.GetRootDirectory(), fileName))
	if err != nil {
		return err
	}
	defer target.Close()

	client := &http.Client{}

	logging.Debug("Downloading: " + url)
	env.DisplayToConsole("Downloading: " + url + "\n")

	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(target, response.Body)
	return err
}
