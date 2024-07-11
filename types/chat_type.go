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
	Sender  string
	Room    string
	Clients map[string]*Client
	Message string
}

type InputNewMessage struct {
	Id     string  `json:"-"`    // Auto
	Room   *string `json:"room"` // Optional and if not filled that indicate create a new room chat and will filled automatically after initialize chat
	UserId string  `json:"-"`    // Auto
	Text   *string `json:"text"`
}

type PushNewChatPayload struct {
	Room         string          `json:"-"`            // Auto
	ChatName     *string         `json:"chat_name"`    // if filled that indicate the chat group otherwise that one to one chat
	Participants []string        `json:"participants"` // is slice of user id
	Message      InputNewMessage `json:"message"`
}
