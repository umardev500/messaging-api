package types

import (
	"sync"
	"time"

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

type BroadcastChatList struct {
	Room      string `json:"room"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type InputNewMessage struct {
	Id     string  `json:"-"`    // Auto
	Room   *string `json:"room"` // Optional and if not filled that indicate create a new room chat and will filled automatically after initialize chat
	UserId string  `json:"-"`    // Auto
	Text   *string `json:"text"`
}

type PushNewChatPayload struct {
	UserId       string          `json:"-"`
	Room         string          `json:"-"`            // Auto
	ChatName     *string         `json:"chat_name"`    // if filled that indicate the chat group otherwise that one to one chat
	Participants []string        `json:"participants"` // is slice of user id
	Message      InputNewMessage `json:"message"`
}

// Chat is chat struct
type Chat struct {
	Id        string     `json:"id"`
	ChatName  *string    `json:"chat_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
