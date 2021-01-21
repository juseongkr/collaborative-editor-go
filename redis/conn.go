package redis

import (
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Connect(address, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
}
