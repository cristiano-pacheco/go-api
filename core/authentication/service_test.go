package authentication_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/cristiano-pacheco/go-api/core/authentication"
	"github.com/cristiano-pacheco/go-api/core/user"
	"github.com/gbrlsnchs/jwt/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

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

func getDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root@/go_api_test?parseTime=true")
	assert.Nil(t, err)
	return db
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

// func TestIssueToken(t *testing.T) {
// 	privKey := "jwt-private-key"
// 	data := newData(1)
// 	db := getDB(t)
// 	defer clearAndClose(db, t)
// 	userService := user.NewService(db, &user.Validator{})
// 	userService.Store(data)
// 	service := authentication.NewService(db, &authentication.Validator{}, privKey)
// 	token, err := service.IssueToken("email1@gmail.com", "password")
// 	assert.Nil(t, err)
// 	assert.IsType(t, &authentication.Token{}, token)
// }

func TestHasAccess(t *testing.T) {
	db := getDB(t)
	service := authentication.NewService(db, &authentication.Validator{}, &jwt.HMACSHA{})
	r, err := service.HasAccess(1, authentication.GetUser)
	assert.True(t, r)
	assert.Nil(t, err)
}
