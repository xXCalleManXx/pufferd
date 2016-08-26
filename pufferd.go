/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"flag"
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/config"
	"github.com/pufferpanel/pufferd/data/templates"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/routing"
	"github.com/pufferpanel/pufferd/routing/legacy"
	"github.com/pufferpanel/pufferd/routing/server"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	var loggingLevel string
	var port int
	flag.StringVar(&loggingLevel, "logging", "DEBUG", "Lowest logging level to display")
	flag.IntVar(&port, "port", 5656, "Port to run service on")
	flag.Parse()

	logging.SetLevelByString(loggingLevel)
	gin.SetMode(gin.ReleaseMode)

	logging.Debug("Logging set to " + loggingLevel)

	config.Load()

	if _, err := os.Stat(templates.Folder); os.IsNotExist(err) {
		logging.Debug("Error on running stat on "+templates.Folder, err)
		err = os.Mkdir(templates.Folder, 755)
		if err != nil {
			logging.Error("Error creating template folder", err)
		}

	}
	if files, _ := ioutil.ReadDir(templates.Folder); len(files) == 0 {
		logging.Debug("Templates being copied to " + templates.Folder)
		templates.CopyTemplates()
	}

	if _, err := os.Stat(programs.ServerFolder); os.IsNotExist(err) {
		logging.Debug("No server directory found, creating", err)
		os.MkdirAll(programs.ServerFolder, 755)
	}

	programs.LoadFromFolder()

	for _, element := range programs.GetAll() {
		if element.IsEnabled() {
			logging.Info("Starting server " + element.Id())
			element.Start()
			err := programs.Save(element.Id())
			if err != nil {
				logging.Error("Error saving server file", err)
			}
		}
	}

	r := gin.Default()
	{
		r.Use(httphandlers.OAuth2Handler)
		routing.RegisterRoutes(r)
		legacy.RegisterRoutes(r)
		server.RegisterRoutes(r)
	}

	var useHttps bool
	useHttps = false

	if _, err := os.Stat(filepath.Join("data", "https.pem")); os.IsNotExist(err) {
		logging.Warn("No HTTPS.PEM found in data folder, will use no http")
	} else if _, err := os.Stat(filepath.Join("data", "https.key")); os.IsNotExist(err) {
		logging.Warn("No HTTPS.KEY found in data folder, will use no http")
	} else {
		useHttps = true
	}

	if useHttps {
		manners.ListenAndServeTLS(":"+strconv.FormatInt(int64(port), 10), filepath.Join("data", "https.pem"), filepath.Join("data", "https.key"), r)
	} else {
		manners.ListenAndServe(":"+strconv.FormatInt(int64(port), 10), r)
	}
}
