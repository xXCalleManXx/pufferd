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

package permissions

import (
	"github.com/pufferpanel/pufferd/utils"
	"regexp"
)

type PermissionTracker interface {
	HasPermission(id string, perm string) bool

	Change(id string, perm string, grant bool)

	Exists(id string) bool

	GetMap() map[string]interface{}
}

type PermTracker struct {
	mapping map[string]interface{}
}

func Create(perms map[string]interface{}) PermissionTracker {
	return &PermTracker{mapping: perms}
}

func (pm *PermTracker) HasPermission(id string, perm string) bool {
	perms := utils.GetStringArrayOrNull(pm.mapping, id)
	if perms == nil {
		return false
	}
	for _, element := range perms {
		if match, _ := regexp.Match(element, []byte(perm)); match {
			return true
		}
	}
	return false
}

func (pm *PermTracker) Change(id string, perm string, grant bool) {
}

func (pm *PermTracker) Exists(id string) bool {
	return utils.GetStringArrayOrNull(pm.mapping, id) != nil
}

func (pm *PermTracker) GetMap() map[string]interface{} {
	return pm.mapping
}
