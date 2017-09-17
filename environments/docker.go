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
	"time"

	ppError "github.com/pufferpanel/pufferd/errors"
	"github.com/docker/docker/client"
	"context"
	"github.com/pufferpanel/apufferi/logging"
)

type docker struct {
	ContainerId   string
	ImageName     string
	BaseEnvironment
}

func (d *docker) ExecuteAsync(cmd string, args []string, callback func(graceful bool)) (error) {
	if d.IsRunning() {
		err := errors.New("A container is already running")
		return err
	}

	return nil
}

func (d *docker) ExecuteInMainProcess(cmd string) (err error) {
	if !d.IsRunning() {
		err = errors.New("main process has not been started")
		return
	}
	//stdIn := d.stdInWriter
	//_, err = io.WriteString(stdIn, cmd+"\n")
	return
}

func (d *docker) Kill() (err error) {
	if !d.IsRunning() {
		return
	}

	return
}

func (d *docker) IsRunning() (bool) {
	client, err := d.getClient()
	if err != nil {
		logging.Error("Error checking run status", err)
		return false
	}
	ctx := context.Background()
	stats, err := client.ContainerInspect(ctx, d.ContainerId)
	if err != nil {
		logging.Error("Error checking run status", err)
		return false
	}
	return stats.State.Running
}

func (d *docker) WaitForMainProcess() error {
	return d.WaitForMainProcessFor(0)
}

func (d *docker) WaitForMainProcessFor(timeout int) (err error) {
	if d.IsRunning() {
		if timeout > 0 {
			var timer = time.AfterFunc(time.Duration(timeout)*time.Millisecond, func() {
				err = d.Kill()
			})
			d.wait.Wait()
			timer.Stop()
		} else {
			d.wait.Wait()
		}
	}
	return
}

func (d *docker) GetStats() (map[string]interface{}, error) {
	if !d.IsRunning() {
		return nil, ppError.NewServerOffline()
	}
	//process, err := process.NewProcess(int32(d.mainProcess.Process.Pid))
	//if err != nil {
	//	return nil, err
	//}
	resultMap := make(map[string]interface{})
	//memMap, _ := process.MemoryInfo()
	resultMap["memory"] = 0
	//cpu, _ := process.Percent(time.Millisecond * 50)
	resultMap["cpu"] = 0
	return resultMap, nil
}

func (d *docker) getClient() (*client.Client, error) {
	return client.NewEnvClient()
}