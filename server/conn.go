package server

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

const (
	writeTimeout   = 10 * time.Second
	readTimeout    = 60 * time.Second
	pingPeriod     = 10 * time.Second
	maxMessageSize = 512
)

type conn struct {
	wsConn   *websocket.Conn
	wg       sync.WaitGroup
	sendChan <-chan interface{}
	id       int
	doc      *server.Document
}

type messageEnvelope struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

func newConn(ws *websocket.Conn, doc *server.Document) *conn {
	return &conn{
		wsConn: ws,
		doc:    doc,
	}
}

func newMessageEnvelop(msg interface{}) (env messageEnvelope) {
	env.Message = msg
	switch msg.(type) {
	case server.InitMessage:
		env.Type = "initMessage"
	case server.OpMessage:
		env.Type = "opMessage"
	case server.ErrorMessage:
		env.Type = "errorMessage"
	default:
		env.Type = "unknownMessage"
	}

	return
}

func (c *conn) run() error {
	id, clientChannel := c.doc.NewClient()
	c.sendChan = clientChannel
	c.id = id

	c.wg.Add(2)
	go c.readPump()
	go c.writePump()
	c.wg.Wait()
	c.doc.RemoveCliend(c.id)

	return nil
}

func (c *conn) readPump() {
	defer c.wg.Done()

	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
	c.wsConn.SetPongHandler(func(string) error {
		c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		var msg server.OpMessage
		if err := c.wsConn.ReadJson(&msg); err != nil {
			log.Println("err reading:", err)
			return
		}

		c.doc.ReceiveChan() <- server.ClienMessage{
			Message:  msg,
			ClientID: c.id,
		}
	}
}

func (c *conn) writePump() {
	defer c.wg.Done()

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case s, more := <-c.sendChan:
			if !more {
				c.wsConn.Close()
				return
			}

			c.wsConn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.wsConn.WriteJSON(newMessageEnvelop(s)); err != nil {
				log.Println("err writing:", err)
				return
			}

		case <-ticker.C:
			c.wsConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout))
		}
	}
}
