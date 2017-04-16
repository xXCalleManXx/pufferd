package sftp

import (
	"github.com/pkg/sftp"
	"path/filepath"
	"github.com/pufferpanel/pufferd/utils"
	"strings"
	"errors"
	"io"
	"os"
	"fmt"
)

type requestPrefix struct {
	prefix string
}

func CreateRequestPrefix(prefix string) sftp.Handlers {
	h := requestPrefix{prefix: prefix}

	return sftp.Handlers{h, h, h, h}
}

func (rp requestPrefix) Fileread(request sftp.Request) (io.ReaderAt, error) {
	file, err := rp.getFile(request.Filepath, os.O_RDONLY, 644)
	return file, err
}

func (rp requestPrefix) Filewrite(request sftp.Request) (io.WriterAt, error) {
	file, err := rp.getFile(request.Filepath, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 644)
	return file, err
}

func (rp requestPrefix) Filecmd(request sftp.Request) error {
	sourceName, err := rp.validate(request.Filepath)
	if err != nil {
		return rp.maskError(err)
	}
	var targetName string
	if request.Target != ""	{
		targetName, err = rp.validate(request.Target)
		if err != nil {
			return rp.maskError(err)
		}
	}
	switch (request.Method) {
	case "SetStat": {
		return nil;
	}
	case "Rename": {
		err = os.Rename(sourceName, targetName)
		return err
	}
	case "Rmdir": {
		return os.RemoveAll(sourceName)
	}
	case "Symlink": {
		return nil;
	}
	default:
		return errors.New(fmt.Sprint("Unknown request method: %s", request.Method));
	}
}

func (rp requestPrefix) Fileinfo(request sftp.Request) ([]os.FileInfo, error) {
	sourceName, err := rp.validate(request.Filepath)
	if err != nil {
		return nil, rp.maskError(err)
	}
	switch (request.Method) {
	case "List": {
		file, err := os.Open(sourceName)
		if err != nil {
			return nil, rp.maskError(err)
		}
		return file.Readdir(0)
	}
	case "Stat": {
		file, err := os.Open(sourceName)
		if err != nil {
			return nil, rp.maskError(err)
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, rp.maskError(err)
		}
		return []os.FileInfo{fi}, nil
	}
	case "Readlink": {
		target, err := os.Readlink(sourceName)
		if err != nil {
			return nil, rp.maskError(err)
		}
		file, err := os.Open(target)
		if err != nil {
			return nil, rp.maskError(err)
		}
		fi, err := file.Stat()
		if err != nil {
			return nil, rp.maskError(err)
		}
		return []os.FileInfo{fi}, nil
	}
	default:
		return nil, errors.New(fmt.Sprint("Unknown request method: %s", request.Method));
	}
}

func (rp requestPrefix) getFile(path string, flags int, mode os.FileMode) (*os.File, error) {
	filePath, err := rp.validate(path)
	if err != nil {
		return nil, rp.maskError(err)
	}
	file, err := os.OpenFile(filePath, flags, mode)
	if err != nil {
		return nil, rp.maskError(err)
	}
	return file, err
}

func (rp requestPrefix) validate(path string) (string, error) {
	ok, path := rp.tryPrefix(path)
	if !ok {
		return "", errors.New("Access denied")
	}
	return path, nil
}

func (rp requestPrefix) tryPrefix(path string) (bool, string) {
	newPath := filepath.Clean(filepath.Join(rp.prefix, path))
	if utils.EnsureAccess(newPath, rp.prefix) {
		return true, newPath
	} else {
		return false, ""
	}
}

func (rp requestPrefix) stripPrefix(path string) string {
	newStr := strings.Replace(path, rp.prefix, "", -1)
	if len(newStr) == 0 {
		newStr = "/"
	}
	return newStr
}

func (rp requestPrefix) maskError(err error) error {
	return errors.New(rp.stripPrefix(err.Error()))
}