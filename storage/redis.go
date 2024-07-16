package storage

import (
	"github.com/redis/go-redis/v9"
)

// var Redis = redis.New(
// 	redis.Config{
// 		Host:      "127.0.0.1",
// 		Port:      6379,
// 		Username:  "",
// 		Password:  "",
// 		Database:  0,
// 		Reset:     false,
// 		TLSConfig: nil,
// 		PoolSize:  10 * runtime.GOMAXPROCS(0),
// 	},
// )

var Redis = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})
