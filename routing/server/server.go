package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/permissions"
	"github.com/pufferpanel/pufferd/programs"
)

func RegisterRoutes(e *gin.Engine) {
	l2 := e.Group("/server", httphandlers.AdminServerAccessHandler, httphandlers.HasServerAccessHandler)
	{
		l2.PUT("/:id", CreateServer)
		l2.DELETE("/:id", DeleteServer)
	}

	l1 := e.Group("/server", httphandlers.UserServerAccessHandler)
	{
		l1.GET("/:id/start", StartServer)
		l1.GET("/:id/stop", StopServer)
		l1.GET("/:id/install", InstallServer)
	}
}

func StartServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.start")

	if !valid {
		return
	}

	existing.Start()
}

func StopServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.stop")

	if !valid {
		return
	}

	existing.Stop()
}

func CreateServer(c *gin.Context) {
	serverId := c.Param("id")
	privKey := c.Query("privkey")
	serverType := c.Query("type")
	data := make(map[string]interface{}, 0)
	data["memory"] = "1024M"

	if !permissions.GetGlobal().HasPermission(privKey, "server.create") {
		c.AbortWithStatus(403)
		return
	}

	existing := programs.GetFromCache(serverId)

	if existing != nil {
		c.AbortWithStatus(409)
		return
	}

	programs.Create(serverId, serverType, data)
}

func DeleteServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.delete")

	if !valid {
		return
	}

	programs.Delete(existing.Id())
}

func InstallServer(c *gin.Context) {
	valid, existing := handleInitialCallServer(c, "server.install")

	if !valid {
		return
	}

	existing.Install()
}

func handleInitialCallServer(c *gin.Context, perm string) (valid bool, program programs.Program) {
	serverId := c.Param("id")
	privKey := c.Query("privkey")

	if !permissions.GetGlobal().HasPermission(privKey, "server.delete") {
		c.AbortWithStatus(403)
		valid = false
		return
	}

	program, _ = programs.Get(serverId)

	if program == nil {
		c.AbortWithStatus(404)
		valid = false
		return
	}

	valid = true
	return
}
