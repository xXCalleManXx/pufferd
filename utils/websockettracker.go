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
	"github.com/gorilla/websocket"
)

type WebSocketManager interface {
	Register(ws *websocket.Conn)

	Write(msg []byte) (n int, e error)
}

type wsManager struct {
	sockets []websocket.Conn
}

func CreateWSManager() WebSocketManager {
	return &wsManager{sockets: make([]websocket.Conn, 0)}
}

func (ws *wsManager) Register(conn *websocket.Conn) {
	ws.sockets = append(ws.sockets, *conn)
}

func (ws *wsManager) Write(msg []byte) (n int, e error) {
	invalid := make([]int, 0)
	for k, v := range ws.sockets {
		err := v.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			invalid = append(invalid, k)
		}
	}
	if len(invalid) > 0 {
		for b := range invalid {
			if len(ws.sockets) == 1 {
				ws.sockets = make([]websocket.Conn, 0)
			} else {
				ws.sockets = append(ws.sockets[:b], ws.sockets[b+1:]...)
			}
		}
	}
	n = len(msg)
	return
}
