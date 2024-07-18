package types

import "time"

type Message struct {
	Id        string     `json:"id"`
	ChatId    string     `json:"chat_id"`
	UserId    string     `json:"user_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type MessagePayload struct {
	Content string `json:"content"`
}

type CreateMessage struct {
	Id      string `json:"-"` // get from auto generated
	ChatId  string `json:"-"` // get from param
	UserId  string `json:"-"` // get from token
	Content string `json:"content"`
}

// Fetch data
type GetMessageType string

const (
	MessageDown GetMessageType = "down"
	MessageUp   GetMessageType = "up"
)

type GetMessageParams struct {
	Type   GetMessageType `json:"type"`
	ChatId string         `json:"-"`
	Date   string         `json:"date"`
}
