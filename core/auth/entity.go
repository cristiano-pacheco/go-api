package auth

import (
	"github.com/gbrlsnchs/jwt/v3"
)

const (
	GetAllUsersAction string = "get_all_users"
	GetUserAction     string = "get_user"
	StoreUserAction   string = "store_user"
	UpdateUserAction  string = "update_user"
	RemoveUserAction  string = "remove_user"
)

// Token struct
type Token struct {
	Token string `json:"token"`
}

// CustomPayload JWT Payload
type CustomPayload struct {
	jwt.Payload
	UserID int64 `json:"user_id"`
}
