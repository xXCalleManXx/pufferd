package httphandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/logging"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logging.Errorf("Error handling route\n%+v", err)
			}
		}()

		c.Next();
	}
}