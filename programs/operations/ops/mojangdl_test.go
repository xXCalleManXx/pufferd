package ops

import (
	"testing"
)

func TestMojangDlOperationFactory_Create(t *testing.T) {
	var factory OperationFactory

	factory = MojangDlOperationFactory{}

	if factory.Key() != "mojangdl" {
		t.Fail()
		return
	}

	version := "release"
	filename := "server.jar"

	createCmd := CreateOperation{
		OperationArgs: make(map[string]interface{}),
		DataMap: make(map[string]interface{}),
	}

	createCmd.OperationArgs["version"] = version
	createCmd.OperationArgs["target"] = filename


	op := factory.Create(createCmd)

	err := op.Run(nil)

	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}