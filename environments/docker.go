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
	"io"
	"io/ioutil"
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

	exists, err := d.doesContainerExist()

	if err != nil {
		return err
	}

	//container does not exist
	if !exists {
		err = d.pullImage(false)

		if err != nil {
			return err
		}

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
			AutoRemove: true,
			NetworkMode: "host",
			Resources: container.Resources{
			},
		}

		networkConfig := &network.NetworkingConfig{
		}
		_, err = client.ContainerCreate(ctx, config, hostConfig, networkConfig, d.ContainerId)
		if err != nil {
			return err
		}
	}

	config := types.ContainerAttachOptions{
		Stdout: true,
		Stderr: true,
		Stream: true,
	}

	response, err := client.ContainerAttach(ctx, d.ContainerId, config)
	if err != nil {
		return err
	}

	go func() {
		defer response.Close()

		wrapper := d.createWrapper()
		io.Copy(wrapper, response.Reader)
	}()

	startOpts := types.ContainerStartOptions{
	}

	err = client.ContainerStart(ctx, d.ContainerId, startOpts)
	if err != nil {
		return err
	}

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

	exists, err := d.doesContainerExist()
	if !exists {
		return false, err
	}

	ctx := context.Background()

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

	resultMap := make(map[string]interface{})
	resultMap["memory"] = 0
	resultMap["cpu"] = 0
	return resultMap, nil
}

func (d *docker) getClient() (*client.Client, error) {
	return client.NewEnvClient()
}

func (d *docker) doesContainerExist() (bool, error) {
	client, err := d.getClient()
	if err != nil {
		return false, err
	}

	ctx := context.Background()

	opts := types.ContainerListOptions{
		Filters: filters.NewArgs(),
	}

	opts.All = true
	opts.Filters.Add("name", d.ContainerId)
	existingContainers, err := client.ContainerList(ctx, opts)
	if len(existingContainers) == 0 {
		return false, err
	} else {
		return true, err
	}
}

func (d *docker) pullImage(force bool) error {
	exists := false

	client, err := d.getClient()
	ctx := context.Background()

	if err != nil {
		return err
	}

	opts := types.ImageListOptions{
		All: true,
		Filters: filters.NewArgs(),
	}
	opts.Filters.Add("reference", d.ImageName)
	images, err := client.ImageList(ctx, opts)

	if err != nil {
		return err
	}

	if len(images) == 1 {
		exists = true
	}

	logging.Debugf("Does image %v exist? %v", d.ImageName, exists)

	if exists && !force {
		return nil
	}

	op := types.ImagePullOptions{
	}

	logging.Debugf("Downloading image %v", d.ImageName)

	r, err := client.ImagePull(ctx, d.ImageName, op)
	defer r.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(ioutil.Discard, r)
	logging.Debugf("Download image %v", d.ImageName)
	return err
}