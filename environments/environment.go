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
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/apufferi/cache"
	"github.com/pufferpanel/pufferd/utils"
	"time"
	"fmt"
	"io"
	"os"
	"github.com/pufferpanel/apufferi/config"
	"sync"
)

type Environment interface {
	//Executes a command within the environment.
	Execute(cmd string, args []string, callback func(graceful bool)) (stdOut []byte, err error)

	//Executes a command within the environment and immediately return
	ExecuteAsync(cmd string, args []string, callback func(graceful bool)) (err error)

	//Sends a string to the StdIn of the main program process
	ExecuteInMainProcess(cmd string) (err error)

	//Kills the main process, but leaves the environment running.
	Kill() (err error)

	//Creates the environment setting needed to run programs.
	Create() (err error)

	//Deletes the environment.
	Delete() (err error)

	Update() (err error)

	IsRunning() (isRunning bool, err error)

	WaitForMainProcess() (err error)

	WaitForMainProcessFor(timeout int) (err error)

	GetRootDirectory() string

	GetConsole() (console []string, epoch int64)

	GetConsoleFrom(time int64) (console []string, epoch int64)

	AddListener(ws *websocket.Conn)

	GetStats() (map[string]interface{}, error)

	DisplayToConsole(msg string, data ...interface{})
}

type BaseEnvironment struct {
	Environment
	RootDirectory string
	ConsoleBuffer cache.Cache
	WSManager utils.WebSocketManager
	wait          sync.WaitGroup
}

func (e *BaseEnvironment) Execute(cmd string, args []string, callback func(graceful bool)) (stdOut []byte, err error) {
	stdOut = make([]byte, 0)
	err = e.ExecuteAsync(cmd, args, callback)
	if err != nil {
		return
	}
	err = e.WaitForMainProcess()
	return
}

func (e *BaseEnvironment) WaitForMainProcess() error {
	return e.WaitForMainProcessFor(0)
}

func (e *BaseEnvironment) WaitForMainProcessFor(timeout int) (err error) {
	running, err := e.IsRunning()
	if err != nil {
		return
	}
	if running {
		if timeout > 0 {
			var timer = time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
				err = e.Kill()
			})
			e.wait.Wait()
			timer.Stop()
		} else {
			e.wait.Wait()
		}
	}
	return
}

func (e *BaseEnvironment) GetRootDirectory() string {
	return e.RootDirectory
}

func (e *BaseEnvironment) GetConsole() (console []string, epoch int64) {
	console, epoch = e.ConsoleBuffer.Read()
	return
}

func (e *BaseEnvironment) GetConsoleFrom(time int64) (console []string, epoch int64) {
	console, epoch = e.ConsoleBuffer.ReadFrom(time)
	return
}

func (e *BaseEnvironment) AddListener(ws *websocket.Conn) {
	e.WSManager.Register(ws)
}

func (e *BaseEnvironment) DisplayToConsole(msg string, data ...interface{}) {
	if len(data) == 0 {
		fmt.Fprint(e.ConsoleBuffer, msg)
		fmt.Fprint(e.WSManager, msg)
	} else {
		fmt.Fprintf(e.ConsoleBuffer, msg, data...)
		fmt.Fprintf(e.WSManager, msg, data...)
	}
}

func (e *BaseEnvironment) Create() error {
	return os.Mkdir(e.RootDirectory, 0755)
}

func (e *BaseEnvironment) Update() error {
	return nil
}

func (e *BaseEnvironment) Delete() (err error) {
	err = os.RemoveAll(e.RootDirectory)
	return
}

func (e *BaseEnvironment) createWrapper() io.Writer {
	if config.Get("forward") == "true" {
		return io.MultiWriter(os.Stdout, e.ConsoleBuffer, e.WSManager)
	}
	return io.MultiWriter(e.ConsoleBuffer, e.WSManager)
}