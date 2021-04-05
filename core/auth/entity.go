package auth

import "github.com/gbrlsnchs/jwt/v3"

// Token struct
type Token struct {
	Token string `json:"token"`
}

// CustomPayload JWT Payload
type CustomPayload struct {
	jwt.Payload
	UserID int64 `json:"user_id"`
}
