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

import "os/exec"

type Environment interface {
	//Starts the environment.
	//This will not start the actual program.
	Start() (err error);

	//Stop the environment.
	//This will not gracefully stop the program.
	//ExecuteInMainProcess should be used before calling this
	Stop() (err error);

	//Executes a command within the environment.
	Execute(cmd string, args ...string) (exitCode int, stdOut []byte, err error);

	//Executes a command within the environment and immediately return
	ExecuteAsync(cmd string, args ...string) (process exec.Cmd, err error);

	//Starts what is considered the main process
	ExecuteMainProcess(cmd string, args ...string) (err error);

	//Sends a string to the StdIn of the main program process
	ExecuteInMainProcess(cmd string) (err error);

	//Kills the main process, but leaves the environment running.
	Kill() (err error);

	//Creates the environment setting needed to run programs.
	Create() (err error);

	//Deletes the environment.
	Delete() (err error);

	//Updates the environment settings.
	//This is similar to recreating the environment without losing data.
	Update() (err error);

	IsRunning() (isRunning bool, err error);
}