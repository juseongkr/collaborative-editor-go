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

func compose(this js.Value, args []js.Value) interface{} {
	ops := make([]ot.CompositeOp, len(args))
	for i, arg := range args {
		if arg.Type() != js.TypeString {
			panic("wrong type")
		}

		var op ot.CompositeOp
		if err := json.Unmarshal([]byte(arg.String()), &op); err != nil {
			panic(err)
		}

		ops[i] = op
	}

	composed := ot.CompositeOp(ops...)
	j, err := json.Marshal(composed)
	if err != nil {
		panic(err)
	}

	return string(j)
}

func (c clientObj) applyOp(op ot.CompositeOp) {
	i, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}

	c.applyOperation.Invoke(string(i))
}

func (c clientObj) sendOp(op ot.CompositeOp, rev int) {
	i, err := json.Marshal(struct {
		Op  ot.CompositeOp `json:"op"`
		Rev int            `json:"rev"`
	}{
		Op:  op,
		Rev: rev,
	})

	if err != nil {
		panic(err)
	}

	c.sendOperation.Invoke(string(i))
}

func newClient(this js.Value, args []js.Value) interface{} {
	id := clientID
	clientID++

	if len(args) != 1 {
		panic("wrong number of arguments")
	}

	opts := args[0]
	if opts.Type() != js.TypeObject {
		panic("wrong argument type")
	}

	rev := opts.Get("revision")
	if rev.Type() != js.TypeNumber {
		panic("revision has to be a number")
	}

	applyOp := opts.Get("applyOperation")
	if applyOp.Type() != js.TypeFunction {
		panic("applyOperation has to be a function")
	}

	sendOp := opts.Get("sendOperation")
	if sendOp.Type() != js.TypeFunction {
		panic("sendOperation has to be a function")
	}

	clients[id] = &clientObject{
		state:          client.State{Revision: rev.Int()},
		applyOperation: applyOp,
		sendOperation:  sendOp,
	}

	return id
}

func getClient(arg js.Value) *clientObj {
	id := arg.Int()
	c, ok := clients[id]
	if !ok {
		panic("no such client")
	}

	return c
}

func getOp(arg js.Value) ot.CompositeOp {
	var op ot.CompositeOp
	if err := json.Unmarshal([]byte(arg.String()), &op); err != nil {
		panic(err)
	}

	return op
}

func applyServer(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		panic("wrong number of arguments")
	}

	c := getClient(args[0])
	op := getOp(args[1])

	newState, docOp := c.state.ApplyServerOp(op)
	c.state = newState
	c.applyOp(docOp)

	return nil
}

func serverAck(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		panic("wrong number of arguments")
	}

	c := getClient(args[0])

	newState, sendAwaiting := c.state.ApplyServerAck()
	c.state = newState

	if sendAwaiting {
		c.sendOp(c.state.Awaiting, c.state.Revision)
	}

	return nil
}

func applyClient(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		panic("wrong number of arguments")
	}

	c := getClient(args[0])
	op := getOp(args[1])

	newState, sendAwaiting := c.state.ApplyClientOp(op)
	c.state = newState

	if sendAwaiting {
		c.sendOp(c.state.Awaiting, c.state.Revision)
	}

	return nil
}

func closeClient(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		panic("wrong number of arguments")
	}

	id := args[0].Int()
	if _, ok := clients[id]; ok {
		delete(clients, id)
	}

	return nil
}

func main() {
	js.Global().Set("compose", js.FuncOf(compose))
	js.Global().Set("newClient", js.FuncOf(newClient))
	js.Global().Set("applyServer", js.FuncOf(applyServer))
	js.Global().Set("serverAck", js.FuncOf(serverAck))
	js.Global().Set("applyClient", js.FuncOf(applyClient))
	js.Global().Set("closeClient", js.FuncOf(closeClient))

	<-make(chan struct{})
}
