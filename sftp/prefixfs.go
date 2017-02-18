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

package sftp

import (
	"os"
	"time"
	"path/filepath"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/pufferpanel/sftp"
	"github.com/pkg/errors"
)

type PrefixFileSystem struct {
	sftp.FileSystem
	prefix string
}

func CreateVirtualFs(prefix string) sftp.FileSystem {
	return PrefixFileSystem{prefix: prefix}
}

func (fs PrefixFileSystem) Stat(path string) (os.FileInfo, error) {
	p, err := fs.validate(path)
	if err != nil {
		return nil, err
	}
	return os.Stat(p)
}

func (fs PrefixFileSystem) Lstat(path string) (os.FileInfo, error) {
	p, err := fs.validate(path)
	if err != nil {
		return nil, err
	}
	return os.Lstat(p)
}

func (fs PrefixFileSystem) Mkdir(path string, mode os.FileMode) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	return os.Mkdir(p, mode)
}

func (fs PrefixFileSystem) Remove(path string) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	return os.Remove(p)
}

func (fs PrefixFileSystem) Symlink(target string, link string) error {
	t, err := fs.validate(target)
	if err != nil {
		return err
	}
	l, err := fs.validate(link)
	if err != nil {
		return err
	}
	return os.Symlink(t, l)
}

func (fs PrefixFileSystem) Readlink(path string) (string, error) {
	p, err := fs.validate(path)
	if err != nil {
		return "", err
	}
	return os.Readlink(p)
}

func (fs PrefixFileSystem) OpenFile(path string, flag int, mode os.FileMode) (*os.File, error) {
	p, err := fs.validate(path)
	if err != nil {
		return nil, err
	}
	return os.OpenFile(p, flag, mode)
}

func (fs PrefixFileSystem) Truncate(path string, size int64) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	return os.Truncate(p, size)
}

func (fs PrefixFileSystem) Chmod(path string, mode os.FileMode) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	return os.Chmod(p, mode)
}

func (fs PrefixFileSystem) Chtimes(path string, aTime, mTime time.Time) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	return os.Chtimes(p, aTime, mTime)
}

func (fs PrefixFileSystem) Chown(path string, uid, gid int) error {
	return errors.New("Chown not supported");
}

func (fs PrefixFileSystem) validate(path string) (string, error) {
	ok, path := fs.tryPrefix(path)
	if !ok {
		return "", errors.New("Invalid path provided")
	}
	return path, nil
}

func (fs PrefixFileSystem) tryPrefix(path string) (bool, string) {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	newPath := filepath.Clean(filepath.Join(fs.prefix, path))
	if utils.EnsureAccess(newPath, fs.prefix) {
		return true, newPath
	} else {
		return false, ""
	}
}