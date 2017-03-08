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
	"strings"
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
	fi, e := os.Stat(p)
	if e != nil {
		e = fs.maskError(e)
	}
	return fi, e
}

func (fs PrefixFileSystem) Lstat(path string) (os.FileInfo, error) {
	p, err := fs.validate(path)
	if err != nil {
		return nil, err
	}
	fi, e := os.Lstat(p)
	if e != nil {
		e = fs.maskError(e)
	}
	return fi, e
}

func (fs PrefixFileSystem) Mkdir(path string, mode os.FileMode) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	e := os.Mkdir(p, mode)
	if e != nil {
		e = fs.maskError(e)
	}
	return e
}

func (fs PrefixFileSystem) Remove(path string) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	e := os.Remove(p)
	if e != nil {
		e = fs.maskError(e)
	}
	return e
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
	e := os.Symlink(t, l)
	if e != nil {
		e = fs.maskError(e)
	}
	return e
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
	e := os.Chmod(p, mode)
	if e != nil {
		e = fs.maskError(e)
	}
	return e
}

func (fs PrefixFileSystem) Chtimes(path string, aTime, mTime time.Time) error {
	p, err := fs.validate(path)
	if err != nil {
		return err
	}
	e := os.Chtimes(p, aTime, mTime)
	if e != nil {
		e = fs.maskError(e)
	}
	return e
}

func (fs PrefixFileSystem) Chown(path string, uid, gid int) error {
	return nil
}

func (fs PrefixFileSystem) Rename(oldPath string, newPath string) error {
	path1, err := fs.validate(oldPath)
	if err != nil {
		return err
	}
	path2, err := fs.validate(newPath)
	if err != nil {
		return err
	}
	e := os.Rename(path1, path2)
	if e != nil {
		e = fs.maskError(e)
	}
	return e
}

func (fs PrefixFileSystem) validate(path string) (string, error) {
	ok, path := fs.tryPrefix(path)
	if !ok {
		return "", errors.New("Invalid path provided")
	}
	return path, nil
}

func (fs PrefixFileSystem) tryPrefix(path string) (bool, string) {
	newPath := filepath.Clean(filepath.Join(fs.prefix, path))
	if utils.EnsureAccess(newPath, fs.prefix) {
		return true, newPath
	} else {
		return false, ""
	}
}

func (fs PrefixFileSystem) maskError(err error) error{
	return errors.New(strings.Replace(err.Error(), fs.prefix, "", -1))
}