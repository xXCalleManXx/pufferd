package httphandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/logging"
	"runtime/debug"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Status(500)
				logging.Errorf("Error handling route\n%+v\n%s", err, debug.Stack())
			}
		}()

		c.Next();
	}
}