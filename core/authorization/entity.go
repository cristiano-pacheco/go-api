package authorization

import "time"

// Permissions
const (
	GetAllUsers string = "get_all_users"
	GetUser     string = "get_user"
	StoreUser   string = "store_user"
	UpdateUser  string = "update_user"
	RemoveUser  string = "remove_user"
)

// Permission struct
type Permission struct {
	ID        int64
	Name      string
	Action    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
