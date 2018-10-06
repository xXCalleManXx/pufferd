/*
 Copyright 2018 Padduck, LLC

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

package commands

import (
	"errors"
	"github.com/pufferpanel/apufferi/logging"
	"os"
	"syscall"
	"time"
)

func Shutdown(pid int) {
	proc, err := os.FindProcess(pid)
	if err != nil || proc == nil {
		if err == nil && proc == nil {
			err = errors.New("no process found")
		}
		logging.Error("Error shutting down pufferd", err)
		return
	}
	err = proc.Signal(syscall.Signal(15))
	if err != nil {
		logging.Error("Error shutting down pufferd", err)
		return
	}

	wait := make(chan error)

	waitForProcess(proc, wait)

	err = <-wait

	if err != nil {
		logging.Error("Error shutting down pufferd", err)
		return
	}
	err = proc.Release()
	if err != nil {
		logging.Error("Error shutting down pufferd", err)
		return
	}
}

func waitForProcess(process *os.Process, c chan error) {
	var err error
	timer := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range timer.C {
			err = process.Signal(syscall.Signal(0))
			if err != nil {
				if err.Error() == "os: process already finished" {
					c <- nil
				} else {
					c <- err
				}

				timer.Stop()
			} else {
			}
		}
	}()
}
