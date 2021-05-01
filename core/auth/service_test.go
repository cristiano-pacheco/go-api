package auth_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/cristiano-pacheco/go-api/core/auth"
	"github.com/cristiano-pacheco/go-api/core/user"
	"github.com/gbrlsnchs/jwt/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestIssueToken(t *testing.T) {
	data := newData(1)
	db := getDB(t)
	defer clearAndClose(db, t)
	userService := user.NewService(db, &user.Validator{})
	userService.Store(data)
	service := getAuthService(db)
	token, err := service.IssueToken("email1@gmail.com", "password")
	assert.Nil(t, err)
	assert.IsType(t, &auth.Token{}, token)
}

func TestHasAccess(t *testing.T) {
	db := getDB(t)
	newAuthData(db)
	service := getAuthService(db)
	r, err := service.HasAccess(1, auth.GetUserAction)
	assert.True(t, r)
	assert.Nil(t, err)
	r, err = service.HasAccess(1, auth.GetAllUsersAction)
	assert.True(t, r)
	assert.Nil(t, err)
	clearAndClose(db, t)
}

func TestGetUserPermissionsById(t *testing.T) {
	db := getDB(t)
	newAuthData(db)
	service := getAuthService(db)
	up, err := service.GetUserPermissionsById(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), up.ID)
	assert.Equal(t, "User", up.Name)
	assert.Equal(t, "Get All Users", up.Permissions[0].Name)
	assert.Equal(t, "get_all_users", up.Permissions[0].Code)
	assert.Equal(t, "Get User", up.Permissions[1].Name)
	assert.Equal(t, "get_user", up.Permissions[1].Code)
	assert.Equal(t, 2, len(up.Permissions))
	clearAndClose(db, t)
}

func getAuthService(db *sql.DB) *auth.Service {
	jwtHash := jwt.NewHS256([]byte("jwt-private-key"))
	return auth.NewService(db, &auth.Validator{}, jwtHash)
}

func newData(id int64) *user.User {
	return &user.User{
		ID:       id,
		Name:     "User Test",
		Email:    fmt.Sprintf("email%d@gmail.com", id),
		Password: "password",
		IsActive: true,
		IsAdmin:  true,
	}
}

func newAuthData(db *sql.DB) {
	stmt, _ := db.Prepare("insert into user (id, name, email, password, is_active, is_admin) values (?, ?, ?, ?, ?, ?)")
	stmt.Exec(1, "User", "user@gmail.com", "123", 1, 1)
	stmt, _ = db.Prepare("insert user_permission (user_id, permission_id) values (?, ?)")
	stmt.Exec(1, 1)
	stmt, _ = db.Prepare("insert user_permission (user_id, permission_id) values (?, ?)")
	stmt.Exec(1, 2)
	defer stmt.Close()
}

func clearAndClose(db *sql.DB, t *testing.T) {
	tx, err := db.Begin()
	assert.Nil(t, err)
	_, err = tx.Exec("delete from user")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	db.Close()
}

func getDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root@/go_api_test?parseTime=true")
	assert.Nil(t, err)
	return db
}
