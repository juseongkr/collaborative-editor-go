package server

import (
	"github.com/juseongkr/collaborative-editor-go/ot"
	"sync"
)

type InitMessage struct {
	Document string `json:"document"`
	Revision int    `json:"revision"`
}

type OpMessage struct {
	AuthorId string         `json:"authorID"`
	Op       ot.CompositeOp `json:"op"`
	Revision int            `json:"revision"`
}

type ClientMessage struct {
	ClientId int
	Message  OpMessage
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type Document struct {
	state       StateStore
	receiveChan chan ClientMessage

	sendChannelsMux sync.RWMutex
	sendChannels    map[int]chan<- interface{}
	channelCounter  int
}
