package template

import (
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/http"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/programs"
	"os"
)

func RegisterRoutes(e *gin.Engine) {
	l := e.Group("_templates")
	{
		l.GET("", httphandlers.OAuth2Handler("node.templates", false), GetTemplates)
		l.GET("/:id", httphandlers.OAuth2Handler("node.templates", false), GetTemplate)
		l.POST("/:id", httphandlers.OAuth2Handler("node.templates.edit", false), EditTemplate)
	}
}

func GetTemplates(c *gin.Context) {
	http.Respond(c).Data(programs.GetPlugins()).Send()
}

func GetTemplate(c *gin.Context) {
	name, exists := c.GetQuery("id")
	if !exists || name == "" {
		http.Respond(c).Fail().Code(400).MessageCode(http.INVALIDREQUEST).Message("no template name provided").Send()
		return
	}
	data, err := programs.GetPlugin(name)
	if err != nil {
		if os.IsNotExist(err) {
			http.Respond(c).Fail().Code(404).MessageCode(http.NOFILE).Message("no template with provided name").Send()
		} else {
			http.Respond(c).Fail().Code(500).MessageCode(http.UNKNOWN).Message("error reading template").Send()
		}
	} else {
		http.Respond(c).Code(200).Data(data).Send()
	}
}

func EditTemplate(c *gin.Context) {
	name, exists := c.GetQuery("id")
	if !exists || name == "" {
		http.Respond(c).Fail().Code(400).MessageCode(http.INVALIDREQUEST).Message("no template name provided").Send()
		return
	}
	data, err := programs.GetPlugin(name)
	if err != nil {
		if os.IsNotExist(err) {
			http.Respond(c).Fail().Code(404).MessageCode(http.NOFILE).Message("no template with provided name").Send()
		} else {
			http.Respond(c).Fail().Code(500).MessageCode(http.UNKNOWN).Message("error reading template").Send()
		}
	} else {
		http.Respond(c).Code(200).Data(data).Send()
	}
}