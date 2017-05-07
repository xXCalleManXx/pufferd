package template

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/http"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/programs"
)

func RegisterRoutes(e *gin.Engine) {
	l := e.Group("_templates")
	{
		l.GET("", httphandlers.OAuth2Handler("node.templates", false), GetTemplates)
		//l.GET("/:id", httphandlers.OAuth2Handler("node.templates", false), GetTemplate)
	}
}

func GetTemplates(c *gin.Context) {
	http.Respond(c).Data(programs.GetPlugins()).Send()
}