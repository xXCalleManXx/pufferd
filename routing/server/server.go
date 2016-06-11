package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/httphandlers"
)

func RegisterRoutes(e *gin.Engine) {
	l := e.Group("/server", httphandlers.ProgramAccessHandler)
	{
		l.GET("/:id/start", StartServer)
		l.GET("/:id/stop", StopServer)
		l.GET("/:id/install", InstallServer)
	}
}

func StartServer(c *gin.Context) {
}

func StopServer(c *gin.Context) {
}

func InstallServer(c *gin.Context) {
}
