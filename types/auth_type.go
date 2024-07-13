package types

type LoginPayload struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthData struct {
	Token string `json:"token"`
}

type UserClaim struct {
	Id       interface{} `json:"id"`
	Username string      `json:"username"`
}
