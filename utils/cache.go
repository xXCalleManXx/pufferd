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

package utils

import (
	"strconv"

	"github.com/pufferpanel/pufferd/config"
	"time"
)

type cache struct {
	buffer   []message
	capacity int
}

type Cache interface {
	Read() (cache []string, epoch int64)

	ReadFrom(startTime int64) (cache []string, epoch int64)

	Write(b []byte) (n int, err error)
}

func CreateCache() *cache {
	capacity, err := strconv.Atoi(config.Get("console-buffer"))
	if err != nil {
		capacity = 50
	}
	return &cache{
		buffer:   make([]message, 0),
		capacity: capacity,
	}
}

func (c *cache) Read() (msg []string, lasttime int64) {
	msg, lasttime = c.ReadFrom(0)
	return
}

func (c *cache) ReadFrom(startTime int64) (msg []string, lasttime int64) {
	result := make([]string, 0)
	endTime := time.Now().Unix()
	for _, v := range c.buffer {
		if v.time > startTime {
			result = append(result, v.msg)
		}
	}
	return result, endTime
}

func (c *cache) Write(b []byte) (n int, err error) {
	if len(c.buffer) == c.capacity {
		c.buffer = c.buffer[1:]
	}
	c.buffer = append(c.buffer, message{msg: string(b), time: time.Now().Unix()})
	n = len(b)
	return
}

type message struct {
	msg  string
	time int64
}
