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

package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/itsjamie/gin-cors"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/pkg/errors"
	"strings"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterRoutes(e *gin.Engine) {
	l := e.Group("/server")
	{
		e.Handle("CONNECT", "/:id/console", func(c *gin.Context) {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
		})
		l.Use(httphandlers.OAuth2Handler)
		l.PUT("/:id", CreateServer)
		l.DELETE("/:id", DeleteServer)
		l.POST("/:id", EditServer)
		l.GET("/:id/start", StartServer)
		l.GET("/:id/stop", StopServer)
		l.POST("/:id/install", InstallServer)
		l.GET("/:id/file/*filename", GetFile)
		l.PUT("/:id/file/*filename", PutFile)
		l.POST("/:id/console", PostConsole)
		l.GET("/:id/stats", GetStats)
		l.POST("/:id/reload", ReloadServer)
	}
	e.GET("/network", httphandlers.OAuth2Handler, NetworkServer)
	e.GET("/server/:id/console", cors.Middleware(cors.Config{
		Origins:     "*",
		Credentials: true,
	}), GetConsole)
}

func StartServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.start", true)

	if !valid {
		c.Status(404)
		return
	}

	existing.Start()
}

func StopServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.stop", true)
	wait := c.Param("wait")
	if wait == "" || (wait != "true" && wait != "false") {
		wait = "true"
	}

	if !valid {
		return
	}

	err := existing.Stop()
	if err != nil {
		c.Error(err)
	}

	if wait == "true" {
		err = existing.GetEnvironment().WaitForMainProcess()
		if err != nil {
			c.Error(err)
		}
	}
}

func CreateServer(c *gin.Context) {
	serverId := c.Param("id")
	data := make(map[string]interface{}, 0)
	err := json.NewDecoder(c.Request.Body).Decode(&data)
	serverType := data["type"].(string)

	if err != nil {
		logging.Error("Error decoding JSON body", err)
		c.AbortWithError(400, err)
		return
	}

	handleInitialCallServer(c, "server.create", false)

	existing := programs.GetFromCache(serverId)

	if existing != nil {
		c.AbortWithStatus(409)
		return
	}

	if !programs.Create(serverId, serverType, data) {
		c.AbortWithStatus(500)
	}
}

func DeleteServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.delete", true)

	if !valid {
		return
	}

	programs.Delete(existing.Id())
}

func InstallServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.install", true)

	if !valid {
		return
	}

	c.Status(200)
	go func() {
		existing.Install()
	}()
}

func EditServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.edit", true)

	if !valid {
		return
	}

	data := make(map[string]interface{}, 0)
	json.NewDecoder(c.Request.Body).Decode(&data)

	c.Status(200)
	existing.Edit(data)
}

func GetFile(c *gin.Context) {

	valid, server := handleInitialCallServer(c, "server.file.get", true)

	if !valid {
		return
	}

	targetPath := c.Param("filename")

	targetFile := utils.JoinPath(server.GetEnvironment().GetRootDirectory(), targetPath)

	if !utils.EnsureAccess(targetFile, server.GetEnvironment().GetRootDirectory()) {
		return
	}

	info, err := os.Stat(targetFile)

	if os.IsNotExist(err) {
		c.Status(404)
		return
	}

	if info.IsDir() {
		files, _ := ioutil.ReadDir(targetFile)
		fileNames := make([]interface{}, 0)
		for _, file := range files {
			type FileDesc struct {
				Name      string    `json:"entry,omitempty"`
				Modified  time.Time `json:"date,omitempty"`
				Size      int64     `json:"size,omitempty"`
				File      bool      `json:"-,omitempty"`
				Directory string    `json:"directory,omitempty"`
				Extension string    `json:"extension,omitempty"`
			}
			fileNames = append(fileNames, &FileDesc{
				Name:      file.Name(),
				Size:      file.Size(),
				File:      !file.IsDir(),
				Modified:  file.ModTime(),
				Extension: filepath.Ext(file.Name()),
				Directory: filepath.Dir(file.Name()),
			})
		}
		c.JSON(200, fileNames)
	} else {
		c.File(targetFile)
	}
}

func PutFile(c *gin.Context) {
	valid, server := handleInitialCallServer(c, "server.file.put", true)

	if !valid {
		return
	}

	targetPath := c.Param("filename")

	if targetPath == "" {
		c.Status(404)
		return
	}

	targetFile := utils.JoinPath(server.GetEnvironment().GetRootDirectory(), targetPath)

	if !utils.EnsureAccess(targetFile, server.GetEnvironment().GetRootDirectory()) {
		return
	}

	file, err := os.Create(targetFile)

	if err != nil {
		logging.Error("Error writing file", err)
		return
	}

	_, err = io.Copy(file, c.Request.Body)

	if err != nil {
		logging.Error("Error writing file", err)
	}
}

func PostConsole(c *gin.Context) {
	valid, program := handleInitialCallServer(c, "server.console.send", true)
	if !valid {
		return
	}
	d, _ := ioutil.ReadAll(c.Request.Body)
	cmd := string(d)
	err := program.Execute(cmd)
	if err != nil {
		c.Error(err)
	} else {
		c.Status(200)
	}
}

func GetConsole(c *gin.Context) {
	httphandlers.ParseToken(c.Query("accessToken"), c)
	valid, program := handleInitialCallServer(c, "server.console", true)
	if !valid {
		return
	}
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logging.Error("Error creating websocket", err)
		c.AbortWithError(500, err)
		return
	}
	console := program.GetEnvironment().GetConsole()
	for _, v := range console {
		conn.WriteMessage(websocket.TextMessage, []byte(v))
	}
	program.GetEnvironment().AddListener(conn)
}

func GetStats(c *gin.Context) {
	valid, server := handleInitialCallServer(c, "server.stats", true)

	if !valid {
		return
	}

	results, err := server.GetEnvironment().GetStats()
	if err != nil {
		result := make(map[string]interface{})
		result["error"] = err.Error()
		c.JSON(200, result)
	} else {
		c.JSON(200, results)
	}
}

func ReloadServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.reload", true)

	if !valid {
		c.Status(404)
		return
	}

	programs.Reload(existing.Id())
}

func NetworkServer(c *gin.Context) {

	scopes, _ := c.Get("scopes")
	valid := false
	for _, v := range scopes.([]string) {
		if v == "pufferadmin"{
			valid = true
		}
	}
	if !valid {
		c.AbortWithStatus(401)
		return
	}

	servers := c.DefaultQuery("ids", "")
	if servers == "" {
		c.AbortWithError(400, errors.New("Server ids required"))
		return
	}
	serverIds := strings.Split(servers, ",")
	result := make(map[string]string)
	for _, v := range serverIds {
		program, _ := programs.Get(v)
		if program == nil {
			continue
		}
		result[program.Id()] = program.GetNetwork()
	}
	c.JSON(200, result)
}

func handleInitialCallServer(c *gin.Context, perm string, requireServer bool) (valid bool, program programs.Program) {
	valid = false

	serverId := c.Param("id")
	targetId, _ := c.Get("server_id")

	if targetId != serverId && targetId != "*" {
		c.AbortWithStatus(401)
		return
	}

	program, _ = programs.Get(serverId)

	if requireServer && program == nil {
		c.AbortWithStatus(404)
		valid = false
		return
	}

	scopes, _ := c.Get("scopes")

	for _, v := range scopes.([]string) {
		if v == perm {
			valid = true
		}
	}

	valid = true

	return
}
