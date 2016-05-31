package routing

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
)

var engine *gin.Engine

func RegisterRootRoute(e *gin.Engine) {
	engine = e
	engine.GET("/", func(c *gin.Context) {
		c.String(200, "pufferd is running")
	})
	engine.GET("_shutdown", Shutdown)
}

func Shutdown(context *gin.Context) {
	for _, element := range programs.GetAll() {
		running, _ := element.IsRunning()
		if running {
			logging.Info("Stopping program " + element.Id())
			element.Stop()
		}
	}
	manners.Close()
}
