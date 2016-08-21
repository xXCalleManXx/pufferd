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

	program, err := programs.Get(serverId)

	if err != nil {
		logging.Error("Error testing permissions", err)
		c.AbortWithStatus(403)
		return
	}

	if program == nil {
		c.AbortWithStatus(403)
		return
	}

	/*if !program.GetPermissionManager().Exists(privKey) {
		c.AbortWithStatus(403)
		return
	}*/
}
