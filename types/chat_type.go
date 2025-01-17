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

type Online struct {
	Conn *websocket.Conn
}

var Onlines = make(map[string]*Online)
var Rooms = make(map[string]map[string]*Client)

type MsgType string

var (
	MsgTypeTyping MsgType = "typing"
	MsgTypeChat   MsgType = "chat"
)

type Broadcast struct {
	Type     MsgType `json:"type"`
	Sender   string  `json:"-"`
	Username string  `json:"username,omitempty"`
	Room     string  `json:"room"`
	Message  string  `json:"message,omitempty"`
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

type CreateNewChatPayload struct {
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

// Get chat list
type GetChatListParam struct {
	UserId string `json:"-"`
	Date   string `json:"date"`
}

type ChatList struct {
	Id          string    `json:"id"`
	ChatName    string    `json:"chat_name"`
	Content     string    `json:"content"`
	LastMsgDate time.Time `json:"last_msg_date"`
}

// Signal
type ChatSignal struct {
	Type     MsgType `json:"type"`
	Username string  `json:"username"`
}
