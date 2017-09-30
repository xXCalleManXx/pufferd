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
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/api/types/network"
)

type docker struct {
	ContainerId   string
	ImageName     string
	BaseEnvironment
}

func (d *docker) ExecuteAsync(cmd string, args []string, callback func(graceful bool)) (error) {
	running, err := d.IsRunning()
	if err != nil {
		return err
	}
	if running {
		return errors.New("container is already running")
	}

	client, err := d.getClient()
	ctx := context.Background()

	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.ExactMatch("container.name", d.ContainerId)

	existingContainers, err := client.ContainerList(ctx, opts)

	//container does not exist
	if len(existingContainers) == 0 {
		cmdSlice := strslice.StrSlice{}

		cmdSlice = append(cmdSlice, cmd)

		for _, v := range args {
			cmdSlice = append(cmdSlice, v)
		}

		config := &container.Config{
			AttachStderr: true,
			AttachStdin: true,
			AttachStdout: true,
			Tty: true,
			StdinOnce: false,
			NetworkDisabled: false,
			Cmd: cmdSlice,
			Image: d.ImageName,
		}

		hostConfig := &container.HostConfig{
			//AutoRemove: true,
			NetworkMode: "host",
			Resources: container.Resources{
				//Memory: int64(mem),
			},
		}

		networkConfig := &network.NetworkingConfig{
		}
		_, err = client.ContainerCreate(ctx, config, hostConfig, networkConfig, d.ContainerId)
		if err != nil {
			return err
		}
	}

	startOpts := types.ContainerStartOptions{
	}

	err = client.ContainerStart(ctx, d.ContainerId, startOpts)
	if err != nil {
		return err
	}

	config := types.ContainerAttachOptions{
		Stdin: true,
	}

	response, err := client.ContainerAttach(ctx, d.ContainerId, config)
	defer response.Close()

	if err != nil {
		return err
	}

	response.Conn.Write([]byte(cmd))
	response.Conn.Write([]byte("\n"))

	return err
}

func (d *docker) ExecuteInMainProcess(cmd string) (err error) {
	running, err := d.IsRunning()
	if err != nil {
		return
	}
	if !running {
		err = errors.New("main process has not been started")
		return
	}
	client, err := d.getClient()
	config := types.ContainerAttachOptions{
		Stdin: true,
	}
	ctx := context.Background()
	response, err := client.ContainerAttach(ctx, d.ContainerId, config)
	defer response.Close()

	response.Conn.Write([]byte(cmd))
	response.Conn.Write([]byte("\n"))
	return
}

func (d *docker) Kill() (err error) {
	running, err := d.IsRunning()
	if err != nil {
		return err
	}

	if !running {
		return
	}

	client, err := d.getClient()
	ctx := context.Background()
	err = client.ContainerKill(ctx, d.ContainerId, "SIGKILL")

	return
}

func (d *docker) IsRunning() (bool, error) {
	client, err := d.getClient()
	if err != nil {
		logging.Error("Error checking run status", err)
		return false, err
	}
	ctx := context.Background()

	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(),
	}
	opts.Filters.ExactMatch("container.name", d.ContainerId)
	existingContainers, err := client.ContainerList(ctx, opts)
	if len(existingContainers) == 0 {
		return false, nil
	}

	stats, err := client.ContainerInspect(ctx, d.ContainerId)
	if err != nil {
		logging.Error("Error checking run status", err)
		return false, err
	}
	return stats.State.Running, nil
}

func (d *docker) WaitForMainProcess() error {
	return d.WaitForMainProcessFor(0)
}

func (d *docker) WaitForMainProcessFor(timeout int) (err error) {
	running, err := d.IsRunning()
	if err != nil {
		return err
	}

	if running {
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
	running, err := d.IsRunning()
	if err != nil {
		return nil, err
	}

	if !running {
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