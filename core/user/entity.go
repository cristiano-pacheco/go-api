package user

import "time"

type User struct {
	ID        int64  `json:id`
	Name      string `json:string`
	Email     string `json:email`
	Password  string `json:password`
	IsAdmin   bool   `json:is_admin`
	CreatedAt time.Time
	UpdatedAt time.Time
}
