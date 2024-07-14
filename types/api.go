package types

type ErrDetail struct {
	Field  string      `json:"field"`
	Filter string      `json:"filter"`
	Detail interface{} `json:"detail"`
}
type Error struct {
	Code    CodeName  `json:"code"`
	Details ErrDetail `json:"details"`
}

type Response struct {
	Code    int         `json:"-"`
	Ticket  string      `json:"ticket"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}
