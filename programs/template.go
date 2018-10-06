package programs

import (
	"fmt"
	"github.com/pkg/errors"
)

type ProgramTemplate struct {
	ProgramData
	SupportedEnvironments []map[string]interface{}
}

func (pt ProgramTemplate) Create(environment string) (Program, error) {
	newPrg := &ProgramData{}
	newPrg.CopyFrom(&pt.ProgramData)
	for _, v := range pt.SupportedEnvironments {
		t := v["type"].(string)
		if t == environment {
			newPrg.EnvironmentData = v
		}
	}

	if newPrg.EnvironmentData == nil {
		return nil, errors.New(fmt.Sprintf("environment specified is not supported (%s)", environment))
	}

	return newPrg, nil
}