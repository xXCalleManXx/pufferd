package httphandlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
)

func UserServerAccessHandler(c *gin.Context) {
	privKey := c.Query("privkey")
	serverId := c.Param("id")

	if serverId == "" {
		c.AbortWithStatus(400)
		return
	}

	if privKey == "" {
		c.Header("WWW-Authentication", "Basic realm=\"pufferd\"")
		c.AbortWithStatus(401)
		return
	}

	program, err := programs.GetProgram(serverId)

	if err != nil {
		logging.Error("Error testing permissions", err)
		c.AbortWithStatus(403)
		return
	}

	if program == nil {
		c.AbortWithStatus(403)
		return
	}

	if !program.GetPermissionManager().Exists(privKey) {
		c.AbortWithStatus(403)
		return
	}
}
