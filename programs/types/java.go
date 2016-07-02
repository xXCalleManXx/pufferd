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

package types

import (
	"encoding/json"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/install"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/permissions"
	"github.com/pufferpanel/pufferd/utils"
	"io/ioutil"
	"os"
)

type Java struct {
	RunData     JavaRun
	InstallData install.InstallSection
	Environment environments.Environment
	Identifier  string
	Permissions permissions.PermissionTracker
	Data        map[string]string
}

//Starts the program.
//This includes starting the environment if it is not running.
func (p *Java) Start() (err error) {
	p.Environment.ExecuteAsync("java", utils.ReplaceTokensInArr(p.RunData.Arguments, p.Data))
	return
}

//Stops the program.
//This will also stop the environment it is ran in.
func (p *Java) Stop() (err error) {
	err = p.Environment.ExecuteInMainProcess(p.RunData.Stop)
	return
}

//Kills the program.
//This will also stop the environment it is ran in.
func (p *Java) Kill() (err error) {
	err = p.Environment.Kill()
	return
}

//Creates any files needed for the program.
//This includes creating the environment.
func (p *Java) Create() (err error) {
	err = p.Environment.Create()
	return
}

//Destroys the server.
//This will delete the server, environment, and any files related to it.
func (p *Java) Destroy() (err error) {
	err = p.Environment.Delete()
	return
}

func (p *Java) Update() (err error) {
	err = p.Install()
	return
}

func (p *Java) Install() (err error) {
	if p.IsRunning() {
		p.Stop()
	}

	os.MkdirAll(p.Environment.GetRootDirectory(), os.ModeDir)

	process := install.GenerateInstallProcess(&p.InstallData, p.Environment, p.Data)
	for process.HasNext() {
		err := process.RunNext()
		if err != nil {
			logging.Error("Error running installer: ", err)
			break
		}
	}
	return
}

//Determines if the server is running.
func (p *Java) IsRunning() (isRunning bool) {
	isRunning = p.Environment.IsRunning()
	return
}

//Sends a command to the process
//If the program supports input, this will send the arguments to that.
func (p *Java) Execute(command string) (err error) {
	err = p.Environment.ExecuteInMainProcess(command)
	return
}

func (p *Java) SetEnabled(isEnabled bool) (err error) {
	p.RunData.Enabled = isEnabled
	return
}

func (p *Java) IsEnabled() (isEnabled bool) {
	isEnabled = p.RunData.Enabled
	return
}

func (p *Java) SetEnvironment(environment environments.Environment) (err error) {
	p.Environment = environment
	return
}

func (p *Java) Id() string {
	return p.Identifier
}

func (p *Java) Name() string {
	return "java"
}

func (p *Java) GetPermissionManager() permissions.PermissionTracker {
	return p.Permissions
}

func (p *Java) GetEnvironment() environments.Environment {
	return p.Environment
}

func (p *Java) SetAutoStart(isAutoStart bool) (err error) {
	p.RunData.AutoStart = isAutoStart
	return
}

func (p *Java) IsAutoStart() (isAutoStart bool) {
	return p.RunData.AutoStart
}

func (p *Java) Save(file string) (err error) {
	result := make(map[string]interface{})
	result["data"] = p.Data
	result["install"] = p.InstallData
	result["permissions"] = p.Permissions.GetMap()
	result["run"] = p.RunData
	result["type"] = "java"

	endResult := make(map[string]interface{})
	endResult["pufferd"] = result

	data, err := json.Marshal(endResult)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 664)
	return
}

type JavaRun struct {
	Stop      string   `json:"stop"`
	Pre       []string `json:"pre,omitempty"`
	Post      []string `json:"post,omitempty"`
	Arguments []string `json:"arguments"`
	Enabled   bool     `json:"enabled"`
	AutoStart bool     `json:"autostart"`
}
