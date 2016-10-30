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
	"errors"
	"os"
	"path/filepath"

	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/taruti/sftpd"
)

type VirtualFS struct {
	sftpd.EmptyFS
	Prefix string
}

type vdir struct {
	d *os.File
}

type vfile struct {
	sftpd.EmptyFile
	f *os.File
}

func (rf vfile) Close() error {
	return rf.f.Close()
}
func (rf vfile) ReadAt(bs []byte, pos int64) (int, error) {
	i, e := rf.f.ReadAt(bs, pos)
	if e != nil {
		logging.Error("Error", e)
	}
	return i, e
}

func (rf vfile) WriteAt(bs []byte, offset int64) (int, error) {
	i, e := rf.f.WriteAt(bs, offset)
	if e != nil {
		logging.Error("Error", e)
	}
	return i, e
}

func (rf vfile) FStat() (*sftpd.Attr, error) {
	fis, e := rf.f.Stat()
	fi := &sftpd.Attr{}
	fi.FillFrom(fis)
	return fi, e
}

func (d vdir) Readdir(count int) ([]sftpd.NamedAttr, error) {
	fis, e := d.d.Readdir(count)
	if e != nil {
		return nil, e
	}
	rs := make([]sftpd.NamedAttr, len(fis))
	for i, fi := range fis {
		rs[i].Name = fi.Name()
		rs[i].FillFrom(fi)
	}
	return rs, nil
}
func (d vdir) Close() error {
	return d.d.Close()
}

func (fs VirtualFS) prefix(path string) (string, error) {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}
	newPath := filepath.Clean(filepath.Join(fs.Prefix, path))
	if utils.EnsureAccess(newPath, fs.Prefix) {
		return newPath, nil
	} else {
		return "<invalid>", errors.New("Invalid path")
	}
}

func (fs VirtualFS) OpenDir(path string) (sftpd.Dir, error) {
	p, e := fs.prefix(path)
	if e != nil {
		return nil, e
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	return vdir{f}, nil
}

func (fs VirtualFS) OpenFile(path string, mode uint32, a *sftpd.Attr) (sftpd.File, error) {
	p, e := fs.prefix(path)
	if e != nil {
		return nil, e
	}
	f, e := os.OpenFile(p, os.O_RDWR, os.ModeType)
	if e != nil {
		if mode == 26 || mode == 58 {
			logging.Debug("Creating file " + p)
			f, e = os.Create(p)
		}
		if e != nil {
			logging.Error("Error openning file", e)
			return nil, e
		}
	}
	return vfile{f: f}, nil
}

func (fs VirtualFS) Stat(name string, islstat bool) (*sftpd.Attr, error) {
	p, e := fs.prefix(name)
	if e != nil {
		return nil, e
	}
	var fi os.FileInfo
	if islstat {
		fi, e = os.Lstat(p)
	} else {
		fi, e = os.Stat(p)
	}
	if e != nil {
		return nil, e
	}
	var a sftpd.Attr
	e = a.FillFrom(fi)

	return &a, e
}

func (fs VirtualFS) SetStat(path string, attr *sftpd.Attr) error {
	return nil
}

func (fs VirtualFS) Remove(name string) error {
	p, e := fs.prefix(name)
	if e != nil {
		return e
	}
	return os.Remove(p)
}

func (fs VirtualFS) Rename(oldName string, newName string, mode uint32) error {
	var p1, p2 string
	var e error
	p1, e = fs.prefix(oldName)
	if e != nil {
		logging.Error("Error renaming file", e)
		return e
	}
	p2, e = fs.prefix(newName)
	if e != nil {
		logging.Error("Error renaming file", e)
		return e
	}
	e = os.Rename(p1, p2)
	if e != nil {
		logging.Error("Error renaming file", e)
	}
	return e
}

func (fs VirtualFS) Mkdir(name string, attr *sftpd.Attr) error {
	p, e := fs.prefix(name)
	if e != nil {
		return e
	}
	return os.Mkdir(p, 0755)
}

func (fs VirtualFS) Rmdir(name string) error {
	p, e := fs.prefix(name)
	if e != nil {
		return e
	}
	return os.Remove(p)
}
