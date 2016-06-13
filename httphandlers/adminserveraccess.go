package httphandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/permissions"
)

func AdminServerAccessHandler(c *gin.Context) {
	privKey := c.Query("privkey")

	if privKey == "" {
		c.Header("WWW-Authentication", "Basic realm=\"pufferd\"")
		c.AbortWithStatus(401)
		return
	}

	if !permissions.GetGlobal().Exists(privKey) {
		c.AbortWithStatus(403)
		return
	}

}
