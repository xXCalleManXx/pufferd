package types

import (
	"github.com/pufferpanel/pufferd/environments"
)

type Java struct {
	Run JavaRun
}

//Starts the program.
//This includes starting the environment if it is not running.
func (p *Java) Start() (err error) {
	return;
}

//Stops the program.
//This will also stop the environment it is ran in.
func (p *Java) Stop() (err error) {
	return;
}

//Kills the program.
//This will also stop the environment it is ran in.
func (p *Java) Kill() (err error) {
	return;
}

//Creates any files needed for the program.
//This includes creating the environment.
func (p *Java) Create() (err error) {
	return;
}

//Destroys the server.
//This will delete the server, environment, and any files related to it.
func (p *Java) Destroy() (err error) {
	return;
}

func (p *Java) Update() (err error) {
	return;
}

func (p *Java) Install() (err error) {
	return;
}

//Determines if the server is running.
func (p *Java) IsRunning() (isRunning bool, err error) {
	return;
}

//Sends a command to the process
//If the program supports input, this will send the arguments to that.
func (p *Java) Execute(command string) (err error) {
	return;
}

func (p *Java) SetEnabled(isEnabled bool) (err error) {
	return;
}

func (p *Java) IsEnabled() (isEnabled bool) {
	return;
}

func (p *Java) GetEnvironment() (environment environments.Environment, err error) {
	return;
}

func (p *Java) SetEnvironment(environment environments.Environment) (err error) {
	return;
}

type JavaRun struct {
	Stop string
	Pre []string
	Post []string
	Arguments string
}