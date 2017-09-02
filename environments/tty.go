// +build !windows

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

package environments

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"github.com/kr/pty"
	"github.com/pufferpanel/apufferi/logging"
)

type tty struct {
	standard
	mainProcess   *exec.Cmd
	stdInWriter   io.Writer
}

func (s *tty) ExecuteAsync(cmd string, args []string) (err error) {
	if s.IsRunning() {
		err = errors.New("A process is already running (" + strconv.Itoa(s.mainProcess.Process.Pid) + ")")
		return
	}
	process := exec.Command(cmd, args...)
	process.Dir = s.RootDirectory
	process.Env = append(os.Environ(), "HOME="+s.RootDirectory)
	if err != nil {
		logging.Error("Error starting process", err)
	}
	wrapper := s.createWrapper()

	if err != nil {
		logging.Error("Error starting process", err)
	}
	s.wait = sync.WaitGroup{}
	s.wait.Add(1)
	process.SysProcAttr = &syscall.SysProcAttr{Setctty: true, Setsid: true}
	s.mainProcess = process
	tty, err := pty.Start(process)
	s.stdInWriter = tty
	go func() {
		io.Copy(wrapper, tty)
		process.Wait()
		s.wait.Done()
	}()
	if err != nil /*&& err.Error() != "exit status 1"*/ {
		logging.Error("Error starting process", err)
	}
	return
}