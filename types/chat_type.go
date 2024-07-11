package types

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

type Broadcast struct {
	Sender  *Client
	Room    string
	Clients map[string]*Client
	Message string
}
