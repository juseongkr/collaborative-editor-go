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
