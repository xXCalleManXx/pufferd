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
	"github.com/pufferpanel/pufferd/data/templates"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/permissions"
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/routing"
	"github.com/pufferpanel/pufferd/routing/legacy"
	"github.com/pufferpanel/pufferd/routing/server"
	"os"
	"strconv"
)

func main() {
	var loggingLevel string
	var port int
	flag.StringVar(&loggingLevel, "logging", "INFO", "Lowest logging level to display")
	flag.IntVar(&port, "port", 5656, "Port to run service on")
	flag.Parse()

	logging.SetLevelByString(loggingLevel)

	if _, err := os.Stat(templates.Folder); os.IsNotExist(err) {
		templates.CopyTemplates()
	}
	if _, err := os.Stat(programs.ServerFolder); os.IsNotExist(err) {
		os.MkdirAll(programs.ServerFolder, os.ModeDir)
	}

	programs.LoadFromFolder()

	for _, element := range programs.GetAll() {
		if element.IsEnabled() {
			logging.Info("Starting server " + element.Id())
			element.Start()
		}
	}

	r := gin.Default()
	{
		routing.RegisterRoutes(r)
		legacy.RegisterRoutes(r)
		server.RegisterRoutes(r)
	}

	permissions.GetGlobal()

	manners.ListenAndServe(":"+strconv.FormatInt(int64(port), 10), r)
}
