package httphandlers

import "github.com/gin-gonic/gin"

func HasServerAccessHandler(c *gin.Context) {
	serverId := c.Param("id")

	if serverId == "" {
		c.AbortWithStatus(404)
		return
	}
}
