package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"

	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type Client struct {
	conn         *ws.Conn
	messageQueue chan Message
	OnMessage    func(Message)
	OnBytes      func([]byte)
	OnClose      func()
	closeChan    chan bool
	Context      context.Context
	connected    bool
	Key          string
	writeMutex   *sync.Mutex
}

func (c *Client) Connect(url string) error {
	if c.connected {
		fmt.Fprintln(os.Stderr, "Already connected")
		return nil
	}
	c.messageQueue = make(chan Message)
	c.closeChan = make(chan bool)
	c.writeMutex = new(sync.Mutex)
	c.Context = context.Background()

	conn, _, err := ws.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	conn.SetCloseHandler(func(code int, text string) error {
		if c.OnClose != nil {
			c.OnClose()
		}
		return nil
	})

	c.conn = conn
	c.listenMessage()
	c.connected = true
	c.Key = uuid.NewString()
	return nil
}

func (c *Client) Send(msgType string, msg interface{}) error {
	var msgStr string
	switch s := msg.(type) {
	case string:
		msgStr = s
	default:
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		msgStr = string(jsonMsg)
	}
	c.messageQueue <- Message{Text: msgStr, Type: msgType}
	return nil
}

// startWriter starts message writer from messageQueue AND starts pinger.
func (c *Client) listenMessage() {
	finish := false
	go func() {
		for {
			if c.conn == nil {
				break
			}

			_, bytes, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Socket read error:", err)
				c.Close()
				return
			}

			var msg Message
			if err := json.Unmarshal(bytes, &msg); err == nil {
				if msg.Type == "DISCONNECT" {
					c.Close()
					continue
				}
				if msg.Type == "PING" {
					if err := c.Send("PONG", ""); err != nil {
						fmt.Fprint(os.Stderr, err.Error())
					}
					continue
				}
				if c.OnMessage != nil {
					go c.OnMessage(msg)
				}
			} else {
				if c.OnBytes != nil {
					go c.OnBytes(bytes)
				}
			}

			if finish {
				break
			}
		}
	}()
	go func() {
		// TODO: Check performance on this. This loop might exhaust CPU.
		for {
			select {
			case msg := <-c.messageQueue:
				if c.conn != nil {
					c.writeMessage(msg)
				}
			case <-c.closeChan:
				finish = true
			}
			if finish {
				break
			}
		}
	}()
}

func (c *Client) writeMessage(msg Message) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	if err := c.conn.WriteJSON(msg); err != nil {
		fmt.Fprintln(os.Stderr, "Socket write message err:", err)
		c.Close()
	}
}

func (c *Client) Close() error {
	c.connected = false
	if c.conn == nil {
		return nil
	}

	if c.OnClose != nil {
		c.OnClose()
	}
	//
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	if err := c.conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, "")); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return nil
		}
		return err
	}
	c.closeChan <- true
	c.conn.Close()
	c.conn = nil

	return nil
}
