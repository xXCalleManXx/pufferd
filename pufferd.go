package main

import (
	"github.com/PufferPanel/pufferd/legacy"
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
