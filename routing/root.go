package routing

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/permissions"
	"github.com/pufferpanel/pufferd/programs"
)

func RegisterRoutes(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.String(200, "pufferd is running")
	})
	e.GET("_shutdown", httphandlers.AdminServerAccessHandler, Shutdown)
}

func Shutdown(c *gin.Context) {
	privKey := c.Query("privkey")

	if !permissions.GetGlobal().HasPermission(privKey, "service.stop") {
		c.AbortWithStatus(403)
		return
	}

	for _, element := range programs.GetAll() {
		running, _ := element.IsRunning()
		if running {
			logging.Info("Stopping program " + element.Id())
			element.Stop()
		}
	}
	manners.Close()
}
