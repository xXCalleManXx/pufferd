package permissions

import (
	"github.com/pufferpanel/pufferd/utils"
)

type PermissionTracker interface {
	HasPermission(id string, perm string) bool

	Change(id string, perm string, grant bool)
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
		if element == perm {
			return true
		}
	}
	return false
}

func (pm *PermTracker) Change(id string, perm string, grant bool) {
}
