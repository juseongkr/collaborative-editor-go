package main

import (
	"github.com/joho/godotenv"
	"github.com/juseongkr/collaborative-editor-go/redis"
	"github.com/juseongkr/collaborative-editor-go/server"
	"log"
	"net/http"
)

func main() {
	redis.Connect(os.Getenv("REDIS_URL"), os.Getenv("REDIS_PW"), 0)
	if err := http.ListenAndServe(":" + os.Getenv("PORT", server.Handler())); err != nil {
		log.Fatalln(err)
	}
}
