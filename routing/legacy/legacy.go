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

package legacy

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(e *gin.Engine) {
	l := e.Group("/legacy")
	{
		l.GET("/server", GetServerInfo)
		l.POST("/server", CreateServer)
		l.PUT("/server", UpdateServerInfo)
		l.DELETE("/server", DeleteServer)

		l.GET("/server/power/:action", ServerPower)
		l.POST("/server/console", ServerConsole)
		l.GET("/server/log/:lines", GetServerLog)

		l.GET("/server/file/:file", GetFile)
		l.PUT("/server/file/:file", UpdateFile)
		l.DELETE("/server/file/:file", DeleteFile)

		l.GET("/server/download/:hash", DownloadFile)

		l.GET("/server/directory/:directory", GetDirectory)

		l.PUT("/server/reinstall", ReinstallServer)
		l.GET("/server/reset-password", ResetPassword)
	}
}

func GetServerInfo(c *gin.Context) {
}

func CreateServer(c *gin.Context) {
}

func UpdateServerInfo(c *gin.Context) {
}

func DeleteServer(c *gin.Context) {
}

func ServerPower(c *gin.Context) {
}

func ServerConsole(c *gin.Context) {
}

func GetServerLog(c *gin.Context) {
}

func GetFile(c *gin.Context) {
}

func UpdateFile(c *gin.Context) {
}

func DeleteFile(c *gin.Context) {
}

func DownloadFile(c *gin.Context) {
}

func GetDirectory(c *gin.Context) {
}

func ReinstallServer(c *gin.Context) {
}

func ResetPassword(c *gin.Context) {
}
