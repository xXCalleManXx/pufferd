package routing

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
)

func RegisterRoutes(e *gin.Engine) {
	e.GET("/", func(c *gin.Context) {
		c.String(200, "pufferd is running")
	})
	e.GET("_shutdown", Shutdown)
}

func Shutdown(c *gin.Context) {
	for _, element := range programs.GetAll() {
		running, _ := element.IsRunning()
		if running {
			logging.Info("Stopping program " + element.Id())
			element.Stop()
		}
	}
	manners.Close()
}
