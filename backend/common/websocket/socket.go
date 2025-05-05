package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lonewolf101101/Architect-betting/backend/common/generator"
	"github.com/lonewolf101101/Architect-betting/backend/common/oapi"
)

type Websocket struct {
	connections map[string]*Connection
	MapMutex    sync.RWMutex
	MsgMutex    sync.RWMutex
	OnConnect   func(r *http.Request, conn *Connection) error
}

// New creates new Websocket instance
func New() *Websocket {
	connections := make(map[string]*Connection)
	return &Websocket{
		connections: connections,
		MapMutex:    sync.RWMutex{},
	}
}

func (ws *Websocket) GetConnection(key string) (*Connection, bool) {
	c, ok := ws.connections[key]
	return c, ok
}

func (ws *Websocket) SendToAll(msgType string, msg interface{}) {
	for _, conn := range ws.connections {
		if err := conn.Send(msgType, msg); err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	}
}

func (ws *Websocket) CloseConnection(key string) {
	conn, ok := ws.connections[key]
	if !ok {
		return
	}

	ws.MapMutex.Lock()
	delete(ws.connections, key)
	ws.MapMutex.Unlock()

	if conn.OnClose != nil {
		conn.OnClose()
	}

	conn.writeMutex.Lock()
	defer conn.writeMutex.Unlock()
	if err := conn.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Printf("%v close err: %v", key, err)
	}

	conn.closeChan <- true
	conn.conn.Close()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *Websocket) Handler(w http.ResponseWriter, r *http.Request) {
	wsb, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		fmt.Fprint(os.Stderr, "upgrade err:", err)
		oapi.ClientError(w, http.StatusBadRequest)
		return
	}

	// Create new connection
	// Assign new key
	k := generator.RandomSimpleString(18)
	ws.MapMutex.Lock()
	ws.connections[k] = newConnection(k, wsb, ws)
	ws.MapMutex.Unlock()
	ws.connections[k].startWriter()
	ws.connections[k].conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("%v Close handler: %v %v", k, code, text)
		ws.CloseConnection(k)
		return nil
	})

	if ws.OnConnect != nil {
		if err := ws.OnConnect(r, ws.connections[k]); err != nil {
			ws.CloseConnection(k)
			oapi.ServerError(w, err)
			return
		}
	}

	go func() {
		for {
			if _, ok := ws.connections[k]; !ok {
				break
			}

			_, bytes, err := ws.connections[k].conn.ReadMessage()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Websocket: reader:", err)
				ws.CloseConnection(k)
				break
			}

			var msg Message
			if err := json.Unmarshal(bytes, &msg); err == nil {
				if msg.Type == "DISCONNECT" {
					ws.CloseConnection(k)
					continue
				}
				if msg.Type == "PONG" {
					if ws.connections[k] != nil {

						ws.connections[k].isPonged = true
						continue
					}
				}
				if ws.connections[k].OnMessage != nil {
					ws.connections[k].OnMessage(msg)
				}
			} else {
				if ws.connections[k].OnBytes != nil {
					ws.connections[k].OnBytes(bytes)
				}
			}
		}
	}()
}
