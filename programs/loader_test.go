package programs_test

import (
	"testing"
	"github.com/pufferpanel/pufferd/programs"
)

func TestLoadServer_Java(t *testing.T) {
	data := []byte("{\"pufferd\":{\"type\":\"java\",\"install\":{\"files\":[\"https://hub.spigotmc.org/BuildTools.jar\"],\"pre\":[],\"post\":[\"java -jar buildtools --rev ${version}\",\"mv spigot*.jar server.jar\"]},\"run\":{\"stop\":\"/stop\",\"pre\":[],\"post\":[],\"arguments\":\"-Xmx${maxmem} -jar server.jar\"}}}");
	var program, err = programs.LoadServerFromData(data);
	if (err != nil || program == nil) {
		if (err != nil) {
			t.Error(err);
		} else {
			t.Error("Program return was nil instead of java");
		}
	}
}

func TestLoadServer_Unknown(t *testing.T) {
	data := []byte("{\"pufferd\": {\"type\": \"badserver\"}}");
	var program, err = programs.LoadServerFromData(data);
	if (err != nil || program != nil) {
		if (err != nil) {
			t.Error(err);
		} else {
			t.Error("Program return was not nil");
		}
	}
}