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

// Package legacy implements the legacy API for compatibility with PufferPanel v0.8.x
package legacy

import (
	"github.com/PufferPanel/pufferd/environments/system"
	"github.com/gin-gonic/gin"
)

func GetServerInfo(c *gin.Context) {
}

func CreateServer(c *gin.Context) {
}

func UpdateServerInfo(c *gin.Context) {
}

func DeleteServer(c *gin.Context) {
}

func ServerPower(c *gin.Context) {
	system.StartServer(c)
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
