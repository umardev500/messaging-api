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
