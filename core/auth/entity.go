package auth

import (
	"github.com/gbrlsnchs/jwt/v3"
)

const (
	GetAllUsersAction     string = "get_all_users"
	GetUserAction         string = "get_user"
	StoreUserAction       string = "store_user"
	UpdateUserAction      string = "update_user"
	RemoveUserAction      string = "remove_user"
	GetAllListsAction     string = "get_all_lists"
	GetListAction         string = "get_list"
	StoreListAction       string = "store_list"
	UpdateListAction      string = "update_list"
	RemoveListAction      string = "remove_list"
	GetAllListItemsAction string = "get_all_list_items"
	GetListItemAction     string = "get_list_item"
	StoreLisItemAction    string = "store_list_item"
	UpdateListItemAction  string = "update_list_item"
	RemoveListItemAction  string = "remove_list_item"
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
