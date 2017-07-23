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

	"github.com/pufferpanel/apufferi/common"
	"github.com/pufferpanel/apufferi/logging"
	"github.com/pufferpanel/pufferd/environments"
	"github.com/pufferpanel/pufferd/programs/operations"
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

type programData struct {
	RunData     Runtime                  `json:"run"`
	InstallData operations.Process       `json:"install"`
	UpdateData  operations.Process       `json:"update"`
	Environment environments.Environment `json:"environment"`
	Identifier  string                   `json:"id"`
	Data        map[string]interface{}   `json:"data"`
}

//Starts the program.
//This includes starting the environment if it is not running.
func (p *programData) Start() (err error) {
	logging.Debugf("Starting server %s", p.Id())
	p.Environment.DisplayToConsole("Starting server\n")
	data := make(map[string]interface{})
	for k, v := range p.Data {
		data[k] = v.(map[string]interface{})["value"]
	}

	err = p.Environment.ExecuteAsync(p.RunData.Program, common.ReplaceTokensInArr(p.RunData.Arguments, data))
	if err != nil {
		p.Environment.DisplayToConsole("Failed to start server\n")
	} else {
		//p.Environment.DisplayToConsole("Server started\n")
	}
	return
}

//Stops the program.
//This will also stop the environment it is ran in.
func (p *programData) Stop() (err error) {
	logging.Debugf("Stopping server %s", p.Id())
	err = p.Environment.ExecuteInMainProcess(p.RunData.Stop)
	if err != nil {
		p.Environment.DisplayToConsole("Failed to stop server\n")
	} else {
		p.Environment.DisplayToConsole("Server stopped\n")
	}
	return
}

//Kills the program.
//This will also stop the environment it is ran in.
func (p *programData) Kill() (err error) {
	logging.Debugf("Killing server %s", p.Id())
	err = p.Environment.Kill()
	if err != nil {
		p.Environment.DisplayToConsole("Failed to kill server\n")
	} else {
		p.Environment.DisplayToConsole("Server killed\n")
	}
	return
}

//Creates any files needed for the program.
//This includes creating the environment.
func (p *programData) Create() (err error) {
	logging.Debugf("Creating server %s", p.Id())
	p.Environment.DisplayToConsole("Allocating server\n")
	err = p.Environment.Create()
	p.Environment.DisplayToConsole("Server allocated\n")
	p.Environment.DisplayToConsole("Ready to be installed\n")
	return
}

//Destroys the server.
//This will delete the server, environment, and any files related to it.
func (p *programData) Destroy() (err error) {
	logging.Debugf("Destroying server %s", p.Id())
	err = p.Environment.Delete()
	return
}

func (p *programData) Update() (err error) {
	logging.Debugf("Updating server %s", p.Id())
	process := operations.GenerateProcess(&p.UpdateData, p.Environment, p.Data)
	err = process.Run()
	if err != nil {
		p.Environment.DisplayToConsole("Error running updater, check daemon logs")
	} else {
		p.Environment.DisplayToConsole("Server updated\n")
	}
	return
}

func (p *programData) Install() (err error) {
	logging.Debugf("Installing server %s", p.Id())
	if p.IsRunning() {
		err = p.Stop()
	}

	if err != nil {
		logging.Error("Error stopping server to install: ", err)
		p.Environment.DisplayToConsole("Error stopping server\n")
		return
	}

	p.Environment.DisplayToConsole("Installing server\n")

	os.MkdirAll(p.Environment.GetRootDirectory(), 0755)

	process := operations.GenerateProcess(&p.InstallData, p.Environment, p.Data)
	err = process.Run()
	if err != nil {
		p.Environment.DisplayToConsole("Error running installer, check daemon logs")
	} else {
		p.Environment.DisplayToConsole("Server installed\n")
	}
	return
}

//Determines if the server is running.
func (p *programData) IsRunning() (isRunning bool) {
	isRunning = p.Environment.IsRunning()
	return
}

//Sends a command to the process
//If the program supports input, this will send the arguments to that.
func (p *programData) Execute(command string) (err error) {
	err = p.Environment.ExecuteInMainProcess(command)
	return
}

func (p *programData) SetEnabled(isEnabled bool) (err error) {
	p.RunData.Enabled = isEnabled
	return
}

func (p *programData) IsEnabled() (isEnabled bool) {
	isEnabled = p.RunData.Enabled
	return
}

func (p *programData) SetEnvironment(environment environments.Environment) (err error) {
	p.Environment = environment
	return
}

func (p *programData) Id() string {
	return p.Identifier
}

func (p *programData) GetEnvironment() environments.Environment {
	return p.Environment
}

func (p *programData) SetAutoStart(isAutoStart bool) (err error) {
	p.RunData.AutoStart = isAutoStart
	return
}

func (p *programData) IsAutoStart() (isAutoStart bool) {
	isAutoStart = p.RunData.AutoStart
	return
}

func (p *programData) Save(file string) (err error) {
	logging.Debugf("Saving server %s", p.Id())

	endResult := make(map[string]interface{})
	endResult["pufferd"] = p

	data, err := json.MarshalIndent(endResult, "", "  ")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(file, data, 0664)
	return
}

func (p *programData) Edit(data map[string]interface{}) (err error) {
	for k, v := range data {
		if v == nil || v == "" {
			delete(p.Data, k)
		}

		var elem map[string]interface{}

		if p.Data[k] == nil {
			elem = make(map[string]interface{})
		} else {
			elem = p.Data[k].(map[string]interface{})
		}
		elem["value"] = v

		p.Data[k] = elem
	}
	err = Save(p.Id())
	return
}

func (p *programData) Reload(data Program) {
	logging.Debugf("Reloading server %s", p.Id())
	replacement := data.(*programData)
	p.Data = replacement.Data
	p.InstallData = replacement.InstallData
	p.RunData = replacement.RunData
}

func (p *programData) GetData() map[string]interface{} {
	return p.Data
}

func (p *programData) GetNetwork() string {
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
	Stop string `json:"stop"`
	//Pre       operations.Process `json:"pre,omitempty"`
	//Post      operations.Process `json:"post,omitempty"`
	Program   string   `json:"program"`
	Arguments []string `json:"arguments"`
	Enabled   bool     `json:"enabled"`
	AutoStart bool     `json:"autostart"`
}
