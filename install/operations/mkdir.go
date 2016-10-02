package operations

import (
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/utils"
	"os"
)

type Mkdir struct {
	TargetFile  string
	Environment environments.Environment
}

func (m *Mkdir) Run() error {
	target := utils.JoinPath(m.Environment.GetRootDirectory(), m.TargetFile)
	return os.MkdirAll(target, 0755)
}
