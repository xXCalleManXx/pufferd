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

package docker

import "os/exec"

type Docker struct {
}

func (d *Docker) Start() (err error) {
	return;
}

func (d *Docker) Stop() (err error) {
	return;
}

func (d *Docker) Execute(cmd string, args ...string) (exitCode int, stdOut []byte, err error) {
	return;
}

func (d *Docker) ExecuteAsync(cmd string, args ...string) (process exec.Cmd, err error) {
	return;
}

func (d *Docker) ExecuteMainProcess(cmd string, args ...string) (err error) {
	return;
}

func (d *Docker) ExecuteInMainProcess(cmd string) (err error) {
	return;
}

func (d *Docker) Kill() (err error) {
	return;
}

func (d *Docker) Create() (err error) {
	return;
}

func (d *Docker) Delete() (err error) {
	return;
}

func (d *Docker) IsRunning() (isRunning bool, err error) {
	return;
}

func (d Docker) Update() (err error) {
	return;
}