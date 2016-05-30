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

package programs_test

import (
	"github.com/pufferpanel/pufferd/programs"
	"testing"
)

func TestLoadProgram_Java(t *testing.T) {
	data := []byte("{\"pufferd\":{\"type\":\"java\",\"install\":{\"files\":[\"https://hub.spigotmc.org/BuildTools.jar\"],\"pre\":[],\"post\":[\"java -jar buildtools --rev ${version}\",\"mv spigot*.jar server.jar\"]},\"run\":{\"stop\":\"/stop\",\"pre\":[],\"post\":[],\"arguments\":\"-Xmx${maxmem} -jar server.jar\"}}}")
	var program, err = programs.LoadProgramFromData("asdfasdf", data)
	if err != nil || program == nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Program return was nil instead of java")
		}
	}
}

func TestLoadProgram_Unknown(t *testing.T) {
	data := []byte("{\"pufferd\": {\"type\": \"badserver\"}}")
	var program, err = programs.LoadProgramFromData("asdfasdf", data)
	if err != nil || program != nil {
		if err != nil {
			t.Error(err)
		} else {
			t.Error("Program return was not nil")
		}
	}
}
