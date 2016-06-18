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
