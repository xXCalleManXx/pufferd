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

	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pufferpanel/apufferi/cache"
	"github.com/pufferpanel/apufferi/config"
	"github.com/pufferpanel/apufferi/logging"
	ppError "github.com/pufferpanel/pufferd/errors"
	"github.com/pufferpanel/pufferd/utils"
	"github.com/shirou/gopsutil/process"
)

type standard struct {
	RootDirectory string
	ConsoleBuffer cache.Cache
	WSManager     utils.WebSocketManager
	mainProcess   *exec.Cmd
	stdInWriter   io.Writer
	wait          sync.WaitGroup
}

func (s *standard) Execute(cmd string, args []string) (stdOut []byte, err error) {
	stdOut = make([]byte, 0)
	err = s.ExecuteAsync(cmd, args)
	if err != nil {
		return
	}
	err = s.WaitForMainProcess()
	return
}

func (s *standard) ExecuteAsync(cmd string, args []string) (err error) {
	if s.IsRunning() {
		err = errors.New("A process is already running (" + strconv.Itoa(s.mainProcess.Process.Pid) + ")")
		return
	}
	s.mainProcess = exec.Command(cmd, args...)
	s.mainProcess.Dir = s.RootDirectory
	s.mainProcess.Env = append(os.Environ(), "HOME="+s.RootDirectory)
	wrapper := s.createWrapper()
	s.mainProcess.Stdout = wrapper
	s.mainProcess.Stderr = wrapper
	pipe, err := s.mainProcess.StdinPipe()
	if err != nil {
		logging.Error("Error starting process", err)
	}
	s.stdInWriter = pipe
	s.wait = sync.WaitGroup{}
	s.wait.Add(1)
	err = s.mainProcess.Start()
	go func() {
		s.mainProcess.Wait()
		s.wait.Done()
	}()
	if err != nil && err.Error() != "exit status 1" {
		logging.Error("Error starting process", err)
	} else {
		logging.Debug("Process started (" + strconv.Itoa(s.mainProcess.Process.Pid) + ")")
	}
	return
}

func (s *standard) ExecuteInMainProcess(cmd string) (err error) {
	if !s.IsRunning() {
		err = errors.New("Main process has not been started")
		return
	}
	stdIn := s.stdInWriter
	_, err = io.WriteString(stdIn, cmd+"\n")
	return
}

func (s *standard) Kill() (err error) {
	if !s.IsRunning() {
		return
	}
	err = s.mainProcess.Process.Kill()
	s.mainProcess.Process.Release()
	s.mainProcess = nil
	return
}

func (s *standard) Create() error {
	return os.Mkdir(s.RootDirectory, 0755)
}

func (s *standard) Update() error {
	return nil
}

func (s *standard) Delete() (err error) {
	err = os.RemoveAll(s.RootDirectory)
	return
}

func (s *standard) IsRunning() (isRunning bool) {
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

func (s *standard) WaitForMainProcess() error {
	return s.WaitForMainProcessFor(0)
}

func (s *standard) WaitForMainProcessFor(timeout int) (err error) {
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

func (s *standard) GetRootDirectory() string {
	return s.RootDirectory
}

func (s *standard) GetConsole() (console []string, epoch int64) {
	console, epoch = s.ConsoleBuffer.Read()
	return
}

func (s *standard) GetConsoleFrom(time int64) (console []string, epoch int64) {
	console, epoch = s.ConsoleBuffer.ReadFrom(time)
	return
}

func (s *standard) AddListener(ws *websocket.Conn) {
	s.WSManager.Register(ws)
}

func (s *standard) GetStats() (map[string]interface{}, error) {
	if !s.IsRunning() {
		return nil, ppError.NewServerOffline()
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

func (s *standard) DisplayToConsole(msg string, data ...interface{}) {
	if len(data) == 0 {
		fmt.Fprint(s.ConsoleBuffer, msg)
		fmt.Fprint(s.WSManager, msg)
	} else {
		fmt.Fprintf(s.ConsoleBuffer, msg, data...)
		fmt.Fprintf(s.WSManager, msg, data...)
	}
}

func (s *standard) createWrapper() io.Writer {
	if config.Get("forward") == "true" {
		return io.MultiWriter(os.Stdout, s.ConsoleBuffer, s.WSManager)
	}
	return io.MultiWriter(s.ConsoleBuffer, s.WSManager)
}
