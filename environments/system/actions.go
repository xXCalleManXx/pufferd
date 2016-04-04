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

package system

import (
	"os/exec"
	"io"
	"errors"
)

type System struct {
	mainProcess *exec.Cmd;
}

func (s *System) Start() (err error) {
	//main system does not need to have any startup
	return;
}

func (s *System) Stop() (err error) {
	s.Kill();
	return;
}

func (s *System) Execute(cmd string, args ...string) (exitCode int, stdOut []byte, err error) {
	var process = exec.Command(cmd, args...);
	stdOut, err = process.Output();
	return;
}

func (s *System) ExecuteAsync(cmd string, args ...string) (process *exec.Cmd, err error) {
	process = exec.Command(cmd, args...);
	process.Start();
	return;
}

func (s *System) ExecuteMainProcess(cmd string, args ...string) (err error) {
	s.mainProcess, err = s.ExecuteAsync(cmd, args...);
	return;
}

func (s *System) ExecuteInMainProcess(cmd string) (err error) {
	if (s.mainProcess == nil) {
		err = errors.New("Main process has not been started");
		return;
	}
	var stdIn, processErr = s.mainProcess.StdinPipe();
	if (processErr != nil) {
		err = processErr;
		return;
	}
	io.WriteString(stdIn, cmd);
	return;
}

func (s *System) Kill() (err error) {
	if (s.mainProcess == nil) {
		return;
	}
	err = s.mainProcess.Process.Kill();
	s.mainProcess.Process.Release();
	s.mainProcess = nil;
	return;
}

func (s *System) Create() (err error) {
	return;
}

func (s *System) Delete() (err error) {
	return;
}

func (s *System) IsRunning() (isRunning bool, err error) {
	isRunning = s.mainProcess.Process != nil;
	return;
}

func (s *System) Update() (err error) {
	return;
}