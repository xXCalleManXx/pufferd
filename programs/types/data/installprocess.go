package data

import (
	"github.com/pufferpanel/pufferd/programs/types/data/operations"
	"runtime"
)

func GenerateInstallProcess(data *InstallSection) InstallProcess {
	var directions []interface{}
	switch runtime.GOOS {
	case "windows":
		directions = data.Windows
	case "mac":
		directions = data.Mac
	default:
		directions = data.Linux
	}
	if directions == nil {
		directions = data.Global
	}
	ops := make([]operations.Operation, 0)
	for _, element := range directions {
		var mapping = element.(map[string]interface{})
		switch mapping["type"] {
		case "command":
			for _, element := range toStringArray(mapping["commands"]) {
				ops = append(ops, &operations.Command{Command: element})
			}
		case "download":
			for _, element := range toStringArray(mapping["files"]) {
				ops = append(ops, &operations.Download{File: element})
			}
		}
	}
	return InstallProcess{processInstructions: ops}
}

type InstallProcess struct {
	processInstructions []operations.Operation
}

func (p *InstallProcess) RunNext() error {
	var op operations.Operation
	op, p.processInstructions = p.processInstructions[0], p.processInstructions[1:]
	op.Run()
	return nil
}

func (p *InstallProcess) HasNext() bool {
	return len(p.processInstructions) != 0 && p.processInstructions[0] != nil
}

func toStringArray(element interface{}) []string {
	switch element.(type) {
	case string:
		return []string{element.(string)}
	case []interface{}:
		var arr = make([]string, 0)
		for _, element := range element.([]interface{}) {
			arr = append(arr, element.(string))
		}
		return arr
	default:
		return []string{}
	}
}
