package storage

import (
	"runtime"

	"github.com/gofiber/storage/redis/v3"
)

var Redis = redis.New(
	redis.Config{
		Host:      "127.0.0.1",
		Port:      6379,
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	},
)
