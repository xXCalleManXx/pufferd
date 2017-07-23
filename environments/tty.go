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
	"time"

	"github.com/gorilla/websocket"
	"github.com/kr/pty"
	"github.com/pufferpanel/apufferi/config"
	ppErrors "github.com/pufferpanel/pufferd/errors"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/shirou/gopsutil/process"
	"fmt"
	"github.com/pufferpanel/apufferi/cache"
)

type tty struct {
	RootDirectory string
	ConsoleBuffer cache.Cache
	WSManager     utils.WebSocketManager
	mainProcess   *exec.Cmd
	stdInWriter   io.Writer
	wait          sync.WaitGroup
}

func (s *tty) Execute(cmd string, args []string) (stdOut []byte, err error) {
	stdOut = make([]byte, 0)
	err = s.ExecuteAsync(cmd, args)
	if err != nil {
		return
	}
	err = s.WaitForMainProcess()
	return
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

func (s *tty) ExecuteInMainProcess(cmd string) (err error) {
	if !s.IsRunning() {
		err = errors.New("Main process has not been started")
		return
	}
	stdIn := s.stdInWriter
	_, err = io.WriteString(stdIn, cmd+"\r")
	return
}

func (s *tty) Kill() (err error) {
	if !s.IsRunning() {
		return
	}
	err = s.mainProcess.Process.Kill()
	s.mainProcess.Process.Release()
	s.mainProcess = nil
	return
}

func (s *tty) Create() error {
	return os.Mkdir(s.RootDirectory, 0755)
}

func (s *tty) Update() error {
	return nil
}

func (s *tty) Delete() error {
	return os.RemoveAll(s.RootDirectory)
}

func (s *tty) IsRunning() (isRunning bool) {
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

func (s *tty) WaitForMainProcess() error {
	return s.WaitForMainProcessFor(0)
}

func (s *tty) WaitForMainProcessFor(timeout int) (err error) {
	if s.IsRunning() {
		if timeout > 0 {
			var timer = time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
				err = s.Kill()
			})
			s.wait.Wait()
			timer.Stop()
		} else {
			s.wait.Wait()
		}
	}
	return
}

func (s *tty) GetRootDirectory() string {
	return s.RootDirectory
}

func (s *tty) GetConsole() ([]string, int64) {
	return s.ConsoleBuffer.Read()
}

func (s *tty) GetConsoleFrom(time int64) ([]string, int64) {
	return s.ConsoleBuffer.ReadFrom(time)
}

func (s *tty) AddListener(ws *websocket.Conn) {
	s.WSManager.Register(ws)
}

func (s *tty) GetStats() (map[string]interface{}, error) {
	if !s.IsRunning() {
		return nil, ppErrors.NewServerOffline()
	}
	process, err := process.NewProcess(int32(s.mainProcess.Process.Pid))
	if err != nil {
		return nil, err
	}
	resultMap := make(map[string]interface{})
	memMap, _ := process.MemoryInfo()
	resultMap["memory"] = memMap.RSS
	cpu, _ := process.Percent(time.Millisecond * 50)
	resultMap["cpu"] = cpu
	return resultMap, nil
}

func (s *tty) DisplayToConsole(msg string, data ...interface{}) {
	if len(data) == 0 {
		fmt.Fprint(s.ConsoleBuffer, msg)
	} else {
		fmt.Fprintf(s.ConsoleBuffer, msg, data...)
	}
}

func (s *tty) createWrapper() io.Writer {
	if config.Get("forward") == "true" {
		return io.MultiWriter(os.Stdout, s.ConsoleBuffer, s.WSManager)
	}
	return io.MultiWriter(s.ConsoleBuffer, s.WSManager)
}
