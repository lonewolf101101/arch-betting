package websocket

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
)

type Message struct {
	Text string
	Raw  interface{}
	Type string
}

type Connection struct {
	Key            string
	conn           *ws.Conn
	messageQueue   chan Message
	OnMessage      func(Message)
	OnBytes        func([]byte)
	OnClose        func()
	closeChan      chan bool
	Context        context.Context
	isPonged       bool
	connectionPool *Websocket
	writeMutex     *sync.Mutex
}

func newConnection(key string, wsc *ws.Conn, websocket *Websocket) *Connection {
	return &Connection{
		conn:           wsc,
		messageQueue:   make(chan Message),
		closeChan:      make(chan bool),
		Key:            key,
		Context:        context.Background(),
		isPonged:       true,
		writeMutex:     new(sync.Mutex),
		connectionPool: websocket,
	}
}

func (conn *Connection) Send(msgType string, msg interface{}) error {
	jsonMessage, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	conn.messageQueue <- Message{Text: string(jsonMessage), Type: msgType}
	return nil
}

func (conn *Connection) SendRaw(msgType string, msg []byte, msgString ...string) error {
	message := Message{Raw: base64.StdEncoding.EncodeToString(msg), Type: msgType}
	if len(msgString) > 0 {
		message.Text = msgString[0]
	}
	conn.messageQueue <- message
	return nil
}

// startWriter starts message writer from messageQueue AND starts pinger.
func (c *Connection) startWriter() {
	finish := false
	go func() {
		// TODO: Check performance on this. This loop might exhaust CPU.
		for {
			select {
			case msg := <-c.messageQueue:
				c.writeMessage(msg)
			case <-c.closeChan:
				finish = true
			}
			if finish {
				break
			}
		}
	}()
	go func() {
		for {
			if !c.isPonged {
				c.connectionPool.CloseConnection(c.Key)
			}
			select {
			case c.messageQueue <- Message{Type: "PING", Text: "Nothing"}:
				c.isPonged = false
			case <-c.closeChan:
				finish = true
			}
			if finish {
				break
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func (c *Connection) writeMessage(msg Message) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	if err := c.conn.WriteJSON(msg); err != nil {
		fmt.Fprintln(os.Stderr, "Socket write message err:", err)
		c.connectionPool.CloseConnection(c.Key)
	}
}
