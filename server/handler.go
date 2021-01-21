package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var document *server.Document

func getDoc() *server.Document {
	if document == null {
		document = server.NewDocument(server.NewMemoryStateStore())
		go document.Run()
	}

	return document
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrader(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		if err := newConn(conn, getDoc()).run(); err != nil {
			log.Println(err)
		}
	})
}
