package operations

import (
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"os"
	"path/filepath"
)

type Move struct {
	SourceFile  string
	TargetFile  string
	Environment environments.Environment
}

func (m *Move) Run() error {
	source := utils.JoinPath(m.Environment.GetRootDirectory(), m.SourceFile)
	target := utils.JoinPath(m.Environment.GetRootDirectory(), m.TargetFile)
	result, valid := validateMove(source, target)
	if !valid {
		return nil
	}
	for k, v := range result {
		logging.Debugf("Moving file from %s to %s", source, target)
		err := os.Rename(k, v)
		if err != nil {
			logging.Error("Error moving file", err)
		}
	}
	return nil
}

func validateMove(source string, target string) (result map[string]string, valid bool) {
	result = make(map[string]string)
	sourceFiles, _ := filepath.Glob(source)
	info, err := os.Stat(target);

	if (err != nil) {
		if os.IsNotExist(err) && len(sourceFiles) > 1 {
			logging.Error("Target folder does not exist", err)
			valid = false
			return
		} else if !os.IsNotExist(err) {
			valid = false
			logging.Error("Error reading target file on move", err)
			return
		}
	} else if info.IsDir() && len(sourceFiles) > 1 {
		logging.Error("Cannot move multiple files to single file target")
		valid = false
		return
	}

	if info != nil && info.IsDir() {
		for _, v := range sourceFiles {
			_, fileName := filepath.Split(v)
			result[v] = filepath.Join(target, fileName)
		}
	} else {
		for _, v := range sourceFiles {
			result[v] = target
		}
	}
	valid = true
	return
}