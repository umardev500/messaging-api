package config

import (
	"sync"

	"github.com/umardev500/messaging-api/utils"
)

var (
	configOnce     sync.Once
	configInstance *AppConfig
)

type Server struct {
	Protocol string
	Port     string
	Host     string
}

type Upload struct {
	Path      string
	StaticUrl string
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type AppConfig struct {
	Server
	Upload
	Database
}

func GetConfig() *AppConfig {
	configOnce.Do(func() {
		var server = Server{
			Protocol: utils.GetEnv("PROTOCOL", "http://"),
			Port:     utils.GetEnv("PORT", "3000"),
			Host:     utils.GetEnv("HOST", "0.0.0.0"),
		}

		configInstance = &AppConfig{
			Server: server,
			Upload: Upload{
				Path:      utils.GetEnv("UPLOAD_PATH", "public/uploads/"),
				StaticUrl: server.Protocol + server.Host + ":" + server.Port + "/api/static/",
			},
		}
	})

	return configInstance
}
