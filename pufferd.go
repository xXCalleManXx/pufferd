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
	"github.com/pufferpanel/pufferd/legacy"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Legacy API for almost drop in compatibility with PufferPanel
	l := r.Group("/legacy")
	{
		l.GET("/", func(c *gin.Context) {
			c.String(200, "pufferd is running")
		})

		l.GET("/server", legacy.GetServerInfo)
		l.POST("/server", legacy.CreateServer)
		l.PUT("/server", legacy.UpdateServerInfo)
		l.DELETE("/server", legacy.DeleteServer)

		l.GET("/server/power/:action", legacy.ServerPower)
		l.POST("/server/console", legacy.ServerConsole)
		l.GET("/server/log/:lines", legacy.GetServerLog)

		l.GET("/server/file/:file", legacy.GetFile)
		l.PUT("/server/file/:file", legacy.UpdateFile)
		l.DELETE("/server/file/:file", legacy.DeleteFile)

		l.GET("/server/download/:hash", legacy.DownloadFile)

		l.GET("/server/directory/:directory", legacy.GetDirectory)

		l.PUT("/server/reinstall", legacy.ReinstallServer)
		l.GET("/server/reset-password", legacy.ResetPassword)

	}

	var port string = ":5656"
	r.Run(port)
}
