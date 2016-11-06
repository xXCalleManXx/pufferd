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

package programs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/install"
	"github.com/pufferpanel/pufferd/logging"
	"github.com/pufferpanel/pufferd/utils"
)

type Program interface {
	//Starts the program.
	//This includes starting the environment if it is not running.
	Start() (err error)

	//Stops the program.
	//This will also stop the environment it is ran in.
	Stop() (err error)

	//Kills the program.
	//This will also stop the environment it is ran in.
	Kill() (err error)

	//Creates any files needed for the program.
	//This includes creating the environment.
	Create() (err error)

	//Destroys the server.
	//This will delete the server, environment, and any files related to it.
	Destroy() (err error)

	Update() (err error)

	Install() (err error)

	//Determines if the server is running.
	IsRunning() (isRunning bool)

	//Sends a command to the process
	//If the program supports input, this will send the arguments to that.
	Execute(command string) (err error)

	SetEnabled(isEnabled bool) (err error)

	IsEnabled() (isEnabled bool)

	SetAutoStart(isAutoStart bool) (err error)

	IsAutoStart() (isAutoStart bool)

	SetEnvironment(environment environments.Environment) (err error)

	Id() string

	GetEnvironment() environments.Environment

	Save(file string) (err error)

	Edit(data map[string]interface{}) (err error)

	Reload(data Program)

	GetData() map[string]interface{}

	GetNetwork() string
}

type ProgramStruct struct {
	RunData     Runtime
	InstallData install.InstallSection
	Environment environments.Environment
	Identifier  string
	Data        map[string]interface{}
}

//Starts the program.
//This includes starting the environment if it is not running.
func (p *ProgramStruct) Start() (err error) {
	logging.Debugf("Starting server %s", p.Id())
	data := make(map[string]interface{})
	for k, v := range p.Data {
		data[k] = v.(map[string]interface{})["value"]
	}
	p.Environment.ExecuteAsync(p.RunData.Program, utils.ReplaceTokensInArr(p.RunData.Arguments, data))
	return
}

//Stops the program.
//This will also stop the environment it is ran in.
func (p *ProgramStruct) Stop() (err error) {
	err = p.Environment.ExecuteInMainProcess(p.RunData.Stop)
	return
}

//Kills the program.
//This will also stop the environment it is ran in.
func (p *ProgramStruct) Kill() (err error) {
	err = p.Environment.Kill()
	return
}

//Creates any files needed for the program.
//This includes creating the environment.
func (p *ProgramStruct) Create() (err error) {
	err = p.Environment.Create()
	return
}

//Destroys the server.
//This will delete the server, environment, and any files related to it.
func (p *ProgramStruct) Destroy() (err error) {
	err = p.Environment.Delete()
	return
}

func (p *ProgramStruct) Update() (err error) {
	err = p.Install()
	return
}

func (p *ProgramStruct) Install() (err error) {
	if p.IsRunning() {
		p.Stop()
	}

	os.MkdirAll(p.Environment.GetRootDirectory(), 0755)

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
func (p *ProgramStruct) IsRunning() (isRunning bool) {
	isRunning = p.Environment.IsRunning()
	return
}

//Sends a command to the process
//If the program supports input, this will send the arguments to that.
func (p *ProgramStruct) Execute(command string) (err error) {
	err = p.Environment.ExecuteInMainProcess(command)
	return
}

func (p *ProgramStruct) SetEnabled(isEnabled bool) (err error) {
	p.RunData.Enabled = isEnabled
	return
}

func (p *ProgramStruct) IsEnabled() (isEnabled bool) {
	isEnabled = p.RunData.Enabled
	return
}

func (p *ProgramStruct) SetEnvironment(environment environments.Environment) (err error) {
	p.Environment = environment
	return
}

func (p *ProgramStruct) Id() string {
	return p.Identifier
}

func (p *ProgramStruct) GetEnvironment() environments.Environment {
	return p.Environment
}

func (p *ProgramStruct) SetAutoStart(isAutoStart bool) (err error) {
	p.RunData.AutoStart = isAutoStart
	return
}

func (p *ProgramStruct) IsAutoStart() (isAutoStart bool) {
	return p.RunData.AutoStart
}

func (p *ProgramStruct) Save(file string) (err error) {
	result := make(map[string]interface{})
	result["data"] = p.Data
	result["install"] = p.InstallData
	result["run"] = p.RunData

	endResult := make(map[string]interface{})
	endResult["pufferd"] = result

	data, err := json.MarshalIndent(endResult, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 0664)
	return
}

func (p *ProgramStruct) Edit(data map[string]interface{}) (err error) {
	for k, v := range data {
		if v == nil || v == "" {
			delete(p.Data, k)
		}
		p.Data[k] = v
	}
	Save(p.Id())
	return
}

func (p *ProgramStruct) Reload(data Program) {
	replacement := data.(*ProgramStruct)
	p.Data = replacement.Data
	p.InstallData = replacement.InstallData
	p.RunData = replacement.RunData
}

func (p *ProgramStruct) GetData() map[string]interface{} {
	return p.Data
}

func (p *ProgramStruct) GetNetwork() string {
	data := p.GetData()
	ip := "0.0.0.0"
	port := "0"

	ipData := data["ip"]
	if ipData != nil {
		ip = ipData.(map[string]interface{})["value"].(string)
	}

	portData := data["port"]
	if portData != nil {
		port = portData.(map[string]interface{})["value"].(string)
	}

	return ip + ":" + port
}

type Runtime struct {
	Stop      string   `json:"stop"`
	Pre       []string `json:"pre,omitempty"`
	Post      []string `json:"post,omitempty"`
	Program   string   `json:"program"`
	Arguments []string `json:"arguments"`
	Enabled   bool     `json:"enabled"`
	AutoStart bool     `json:"autostart"`
}
