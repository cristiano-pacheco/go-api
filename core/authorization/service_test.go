package authorization_test

import (
	"database/sql"
	"testing"

	"github.com/cristiano-pacheco/go-api/core/authorization"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func getDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "root:root@/go_api_test?parseTime=true")
	assert.Nil(t, err)
	return db
}

func newData(db *sql.DB) {
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

func TestHasAccess(t *testing.T) {
	db := getDB(t)
	newData(db)
	service := authorization.NewService(db)
	r, err := service.HasAccess(1, authorization.GetUser)
	assert.True(t, r)
	assert.Nil(t, err)
	r, err = service.HasAccess(1, authorization.GetAllUsers)
	assert.True(t, r)
	assert.Nil(t, err)
	clearAndClose(db, t)
}
