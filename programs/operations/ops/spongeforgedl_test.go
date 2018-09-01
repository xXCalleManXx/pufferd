package ops

import (
	"github.com/pufferpanel/pufferd/cache"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/utils"
	"testing"
)

func TestSpongeForgeDlOperationFactory_Create(t *testing.T) {
	var factory OperationFactory

	factory = SpongeForgeDlOperationFactory{}

	if factory.Key() != "spongeforgedl" {
		t.Fail()
		return
	}

	version := "stable"

	createCmd := CreateOperation{
		OperationArgs: make(map[string]interface{}),
		DataMap:       make(map[string]interface{}),
	}

	createCmd.OperationArgs["releaseType"] = version

	op := factory.Create(createCmd)

	err := op.Run(&environments.BaseEnvironment{
		ConsoleBuffer: cache.CreateCache(),
		WSManager:     utils.CreateWSManager(),
	})

	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}
