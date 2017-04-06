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
	gohttp "net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/itsjamie/gin-cors"
	"github.com/pufferpanel/pufferd/http"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/pkg/errors"
	"strings"
	"strconv"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *gohttp.Request) bool {
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
		l.DELETE("/:id/file/*filename", DeleteFile)
		l.POST("/:id/console", PostConsole)
		l.GET("/:id/stats", GetStats)
		l.POST("/:id/reload", ReloadServer)
		l.GET("/:id/console", cors.Middleware(cors.Config{
			Origins:     "*",
			Credentials: true,
		}), GetConsole)
		l.GET("/:id/logs", GetLogs)
	}
	e.GET("/network", httphandlers.OAuth2Handler, NetworkServer)
}

func StartServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.start", true)

	if !valid {
		rejectConnection(c, "server.start", existing)
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
		rejectConnection(c, "server.stop", existing)
		return
	}

	err := existing.Stop()
	if err != nil {
		errorConnection(c, err)
		return
	}

	if wait == "true" {
		err = existing.GetEnvironment().WaitForMainProcess()
		if err != nil {
			errorConnection(c, err)
			return
		}
	}
	http.Respond(c).Send()
}

func CreateServer(c *gin.Context) {
	serverId := c.Param("id")
	valid, existing := handleInitialCallServer(c, "server.create", false)

	if !valid {
		rejectConnection(c, "server.create", existing)
		return
	}

	if existing != nil {
		http.Respond(c).Code(409).Message("server already exists").Send()
		return
	}

	data := make(map[string]interface{}, 0)
	err := json.NewDecoder(c.Request.Body).Decode(&data)

	if err != nil {
		logging.Error("Error decoding JSON body", err)
		http.Respond(c).Code(400).Message("error parsing json").Data(err).MessageCode(http.MALFORMEDJSON).Send()
		return
	}

	serverType := data["type"].(string)

	if !programs.Create(serverId, serverType, data) {
		errorConnection(c, nil)
	} else {
		http.Respond(c).Send()
	}
}

func DeleteServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.delete", true)

	if !valid {
		rejectConnection(c, "server.delete", existing)
		return
	}

	programs.Delete(existing.Id())
	http.Respond(c).Send()
}

func InstallServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.install", true)

	if !valid {
		rejectConnection(c, "server.instal", existing)
		return
	}

	http.Respond(c).Send()
	go func() {
		existing.Install()
	}()
}

func EditServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.edit", true)

	if !valid {
		rejectConnection(c, "server.edit", existing)
		return
	}

	data := make(map[string]interface{}, 0)
	json.NewDecoder(c.Request.Body).Decode(&data)

	existing.Edit(data)
	http.Respond(c).Send()
}

func GetFile(c *gin.Context) {

	valid, server := handleInitialCallServer(c, "server.file.get", true)

	if !valid {
		rejectConnection(c, "server.file.get", server)
		return
	}

	targetPath := c.Param("filename")

	targetFile := utils.JoinPath(server.GetEnvironment().GetRootDirectory(), targetPath)

	if !utils.EnsureAccess(targetFile, server.GetEnvironment().GetRootDirectory()) {
		http.Respond(c).Code(403).Message("invalid file path").Code(http.NOTAUTHORIZED).Send()
		return
	}

	info, err := os.Stat(targetFile)

	if os.IsNotExist(err) {
		errorConnection(c, err)
		return
	}

	if info.IsDir() {
		type FileDesc struct {
			Name      string    `json:"name"`
			Modified  int64     `json:"modifyTime"`
			Size      int64     `json:"size,omitempty"`
			File      bool      `json:"isFile"`
			Extension string    `json:"extension,omitempty"`
		}

		files, _ := ioutil.ReadDir(targetFile)
		fileNames := make([]interface{}, 0)
		if targetPath != "" && targetPath != "." && targetPath != "/" {
			newFile := &FileDesc{
				Name:      "..",
				File:      false,
			}
			fileNames = append(fileNames, newFile)
		}
		for _, file := range files {
			newFile := &FileDesc{
				Name:      file.Name(),
				File:      !file.IsDir(),
			}

			if newFile.File {
				newFile.Size = file.Size()
				newFile.Modified = file.ModTime().Unix()
				newFile.Extension = filepath.Ext(file.Name())
			}

			fileNames = append(fileNames, newFile)
		}
		http.Respond(c).Data(fileNames).Send()
	} else {
		_, err := os.Open(targetFile)
		if err != nil {
			if err == os.ErrNotExist {
				http.Respond(c).Code(404).MessageCode(http.NOFILE).Send()
			} else {
				errorConnection(c, err)
			}
		}
		c.File(targetFile)
	}
}

func PutFile(c *gin.Context) {
	valid, server := handleInitialCallServer(c, "server.file.put", true)

	if !valid {
		rejectConnection(c, "server.file.put", server)
		return
	}

	targetPath := c.Param("filename")

	if targetPath == "" {
		c.Status(404)
		return
	}

	targetFile := utils.JoinPath(server.GetEnvironment().GetRootDirectory(), targetPath)

	if !utils.EnsureAccess(targetFile, server.GetEnvironment().GetRootDirectory()) {
		http.Respond(c).Code(403).Message("invalid file path").Code(http.NOTAUTHORIZED).Send()
		return
	}

	_, mkFolder := c.GetQuery("folder")
	if (mkFolder) {
		err := os.Mkdir(targetFile, 0644)
		if err != nil {
			errorConnection(c, err)
		} else {
			http.Respond(c).Send()
		}
		return
	}
	file, err := os.Create(targetFile)

	if err != nil {
		errorConnection(c, err)
		logging.Error("Error writing file", err)
		return
	}

	_, noform := c.GetQuery("noform")

	var sourceFile io.ReadCloser

	if noform {
		c.Request.ParseMultipartForm(32 << 20)
		sourceFile, _, err = c.Request.FormFile("file")
	} else {
		sourceFile = c.Request.Body
	}

	_, err = io.Copy(file, sourceFile)

	if err != nil {
		errorConnection(c, err)
		logging.Error("Error writing file", err)
	} else {
		http.Respond(c).Send()
	}
}

func DeleteFile (c *gin.Context) {
	valid, server := handleInitialCallServer(c, "server.file.delete", true)

	if !valid {
		rejectConnection(c, "server.file.delete", server)
		return
	}


	targetPath := c.Param("filename")

	targetFile := utils.JoinPath(server.GetEnvironment().GetRootDirectory(), targetPath)

	if !utils.EnsureAccess(targetFile, server.GetEnvironment().GetRootDirectory()) {
		http.Respond(c).Code(403).Message("invalid file path").Code(http.NOTAUTHORIZED).Send()
		return
	}

	err := os.Remove(targetFile)
	if err != nil {
		errorConnection(c, err)
		logging.Error("Failed to delete file", err)
	} else {
		http.Respond(c).Send()
	}
}

func PostConsole(c *gin.Context) {
	valid, program := handleInitialCallServer(c, "server.console.send", true)
	if !valid {
		rejectConnection(c, "server.console.send", program)
		return
	}

	d, _ := ioutil.ReadAll(c.Request.Body)
	cmd := string(d)
	err := program.Execute(cmd)
	if err != nil {
		errorConnection(c, err)
	} else {
		http.Respond(c).Send()
	}
}

func GetConsole(c *gin.Context) {
	valid, program := handleInitialCallServer(c, "server.console", true)
	if !valid {
		rejectConnection(c, "server.console", program)
		return
	}

	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logging.Error("Error creating websocket", err)
		errorConnection(c, err)
		return
	}
	console, _ := program.GetEnvironment().GetConsole()
	for _, v := range console {
		conn.WriteMessage(websocket.TextMessage, []byte(v))
	}
	program.GetEnvironment().AddListener(conn)
}

func GetStats(c *gin.Context) {
	valid, server := handleInitialCallServer(c, "server.stats", true)

	if !valid {
		rejectConnection(c, "server.stats", server)
		return
	}

	results, err := server.GetEnvironment().GetStats()
	if err != nil {
		result := make(map[string]interface{})
		result["error"] = err.Error()
		http.Respond(c).Data(result).Code(500).Send()
	} else {
		http.Respond(c).Data(results).Send()
	}
}

func ReloadServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.reload", true)

	if !valid {
		rejectConnection(c, "server.reload", existing)
		return
	}

	err := programs.Reload(existing.Id())
	if err != nil {
		errorConnection(c, err)
		return
	}
	http.Respond(c).Send()
}

func NetworkServer(c *gin.Context) {

	scopes, _ := c.Get("scopes")
	valid := false
	for _, v := range scopes.([]string) {
		if v == "server.network"{
			valid = true
		}
	}
	if !valid {
		http.Respond(c).Code(403).MessageCode(http.NOTAUTHORIZED).Message("missing scope server.network").Send()
		return
	}

	servers := c.DefaultQuery("ids", "")
	if servers == "" {
		http.Respond(c).Code(400).MessageCode(http.NOSERVERID).Message("no server ids provided").Send()
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
	http.Respond(c).Data(result).Send()
}

func GetLogs (c *gin.Context) {
	valid, program := handleInitialCallServer(c, "server.console", true)
	if !valid {
		rejectConnection(c, "server.console", program)
		return
	}

	time := c.DefaultQuery("time", "0")

	castedTime, ok := strconv.ParseInt(time, 10, 64)

	if ok != nil {
		c.AbortWithError(400, errors.New("Time provided is not a valid UNIX time"))
		return
	}

	console, epoch := program.GetEnvironment().GetConsoleFrom(castedTime)
	msg := ""
	for _, k := range console {
		msg += k
	}
	result := make(map[string]interface{})
	result["epoch"] = epoch
	result["logs"] = msg;
	http.Respond(c).Data(result).Send()
}

func handleInitialCallServer(c *gin.Context, perm string, requireServer bool) (valid bool, program programs.Program) {
	valid = false

	serverId := c.Param("id")
	canAccessId, _ := c.Get("server_id")

	accessId := canAccessId.(string)

	if accessId == "*" {
		program, _ = programs.Get(serverId)
	} else {
		program, _ = programs.Get(accessId)
	}

	if accessId != serverId && accessId != "*" {
		return
	}

	if requireServer && program == nil {
		return
	}

	scopes, _ := c.Get("scopes")

	for _, v := range scopes.([]string) {
		if v == perm {
			valid = true
		}
	}
	return
}

func rejectConnection(c *gin.Context, scope string, server programs.Program) {
	builder := http.Respond(c)
	if server != nil {
		builder.Code(403).MessageCode(http.NOTAUTHORIZED).Message("missing scope " + scope)
	} else {
		builder.Code(404).MessageCode(http.NOSERVER).Message("no server with id " + c.Param("id"))
	}
	builder.Send()
}

func errorConnection(c *gin.Context, err error) {
	http.Respond(c).Code(500).MessageCode(http.UNKNOWN).Data(err).Message("error handling request").Send()
}