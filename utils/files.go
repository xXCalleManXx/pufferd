package utils

import (
	"path/filepath"
	"strings"
)

func JoinPath(paths ...string) string {
	return strings.Join(paths, string(filepath.Separator))
}
