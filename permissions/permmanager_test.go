package permissions

import (
	"testing"
)

func TestPermTracker_Change(t *testing.T) {

}

func TestPermTracker_HasPermission_UserExistButNotPerm(t *testing.T) {
	permTracker := createTestPermTracker()
	if permTracker.HasPermission("userA", "permissionC") {
		t.Error("userA should not have permissionC")
	}
}

func TestPermTracker_HasPermission_UserExistAndHasPerm(t *testing.T) {
	permTracker := createTestPermTracker()
	if !permTracker.HasPermission("userA", "permissionA") {
		t.Error("userA should have permissionA")
	}
}

func TestPermTracker_HasPermission_UserDoesNotExist(t *testing.T) {
	permTracker := createTestPermTracker()
	if permTracker.HasPermission("userB", "permissionC") {
		t.Error("userB should not have permissionC")
	}
}

func TestPermTracker_Exists_UserExistAndHasPerm(t *testing.T) {
	permTracker := createTestPermTracker()
	if !permTracker.Exists("userA") {
		t.Error("userA should exist")
	}
}

func TestPermTracker_Exists_UserDoesNotExist(t *testing.T) {
	permTracker := createTestPermTracker()
	if permTracker.Exists("userB") {
		t.Error("userB should not exist")
	}
}

func createTestPermTracker() PermissionTracker {
	perms := make(map[string]interface{})
	perms["userA"] = []interface{}{
		"permissionA",
		"permissionB",
	}
	return &PermTracker{mapping: perms}
}
