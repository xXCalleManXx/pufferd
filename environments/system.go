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
	"github.com/pufferpanel/pufferd/logging"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

type System struct {
	mainProcess   *exec.Cmd
	RootDirectory string
}

func (s *System) Execute(cmd string, args []string) (stdOut []byte, err error) {
	if s.IsRunning() {
		err = errors.New("A process is already running (" + strconv.Itoa(s.mainProcess.Process.Pid) + ")")
		return
	}
	s.mainProcess = exec.Command(cmd, args...)
	s.mainProcess.Dir = s.RootDirectory
	s.mainProcess.Stdout = os.Stdout
	s.mainProcess.Stderr = os.Stderr
	err = s.mainProcess.Run()
	go func() {
		s.mainProcess.Wait()
	}()
	if err != nil && err.Error() != "exit status 1" {
		logging.Error("Error starting process", err)
	}
	return
}

func (s *System) ExecuteAsync(cmd string, args []string) (err error) {
	if s.IsRunning() {
		err = errors.New("A process is already running (" + strconv.Itoa(s.mainProcess.Process.Pid) + ")")
		return
	}
	s.mainProcess = exec.Command(cmd, args...)
	s.mainProcess.Dir = s.RootDirectory
	s.mainProcess.Stdout = os.Stdout
	s.mainProcess.Stderr = os.Stderr
	err = s.mainProcess.Start()
	go func() {
		s.mainProcess.Wait()
	}()
	return
}

func (s *System) ExecuteInMainProcess(cmd string) (err error) {
	if !s.IsRunning() {
		err = errors.New("Main process has not been started")
		return
	}
	var stdIn, processErr = s.mainProcess.StdinPipe()
	if processErr != nil {
		err = processErr
		return
	}
	io.WriteString(stdIn, cmd)
	return
}

func (s *System) Kill() (err error) {
	if !s.IsRunning() {
		return
	}
	err = s.mainProcess.Process.Kill()
	s.mainProcess.Process.Release()
	s.mainProcess = nil
	return
}

func (s *System) Create() (err error) {
	os.Mkdir(s.RootDirectory, os.ModeDir)
	return
}

func (s *System) Delete() (err error) {
	return
}

func (s *System) IsRunning() (isRunning bool) {
	isRunning = s.mainProcess != nil && s.mainProcess.Process != nil
	if isRunning {
		process, pErr := os.FindProcess(s.mainProcess.Process.Pid)
		if process == nil || pErr != nil {
			isRunning = false
		} else if process.Signal(syscall.Signal(0)) != nil {
			isRunning = false
		}
	}
	return
}

func (s *System) WaitForMainProcess() (err error) {
	return s.WaitForMainProcessFor(0)
}

func (s *System) WaitForMainProcessFor(timeout int) (err error) {
	if s.IsRunning() {
		if timeout > 0 {
			var timer = time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
				err = s.Kill()
			})
			err = s.mainProcess.Wait()
			timer.Stop()
		} else {
			err = s.mainProcess.Wait()
		}
	}
	return
}

func (s *System) GetRootDirectory() string {
	return s.RootDirectory
}
