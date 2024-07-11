package types

type Error struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
}

type Response struct {
	Code    int         `json:"-"`
	Ticket  string      `json:"ticket"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}
