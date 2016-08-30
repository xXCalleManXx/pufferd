package utils

import (
	"github.com/gorilla/websocket"
	"github.com/pufferpanel/pufferd/logging"
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
	logging.Debug("asdfadsf")
	ws.sockets = append(ws.sockets, *conn)
	logging.Debugf("%s", len(ws.sockets))
}

func (ws *wsManager) Write(msg []byte) (n int, e error) {
	invalid := make([]int, 0)
	for k, v := range ws.sockets {
		err := v.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			logging.Debug("Client detached")
			invalid = append(invalid, k)
		}
	}
	if len(invalid) > 0 {
		for b := range invalid {
			logging.Debugf("Removing client %s", b)
			ws.sockets = append(ws.sockets[:b], ws.sockets[b+1:]...)
		}
	}
	n = len(msg)
	return
}
