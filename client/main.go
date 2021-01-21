package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
)

type clientObj struct {
	state          client.State
	applyOperation js.Value
	sendOperation  js.Value
}

var clientID int
var clients = make(map[int]*clientObj)
