package types

import "time"

type Participant struct {
	ChatId    string     `json:"chat_id"`
	UserId    string     `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
