/*
 Copyright 2016 Padduck, LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 	http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package routing

import (
	"github.com/braintree/manners"
	"github.com/gin-gonic/gin"
	"github.com/pufferpanel/pufferd/httphandlers"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/programs"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/pufferpanel/pufferd/http"
)

func RegisterRoutes(e *gin.Engine) {
	e.Use(httphandlers.Recovery())
	e.GET("", func(c *gin.Context) {
		http.Respond(c).Message("pufferd is running").Send()
	})
	e.GET("/templates", GetTemplates)
	e.GET("/_shutdown", httphandlers.OAuth2Handler, Shutdown)
}

func Shutdown(c *gin.Context) {
	if !hasScope(c, "node.stop") {
		http.Respond(c).Fail().Code(403).MessageCode(http.NOTAUTHORIZED).Message("missing scope node.stop").Send()
		return
	}

	for _, element := range programs.GetAll() {
		running := element.IsRunning()
		if running {
			logging.Warn("Stopping program " + element.Id())
			element.Stop()
		}
	}
	manners.Close()
}

func GetTemplates(c *gin.Context) {
	http.Respond(c).Data(programs.GetPlugins()).Send()
}

func hasScope(gin *gin.Context, scope string) bool {
	scopes, _ := gin.Get("scopes")
	return utils.ContainsValue(scopes.([]string), scope)
}
