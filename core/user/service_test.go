package user_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/cristiano-pacheco/go-api/core/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// var dsn string

// func init() {
// 	flag.StringVar(&dsn, "dsn", "root:root@/go_api_test?parseTime=true", "MySQL data source name")
// 	flag.Parse()
// }

func newData(id int64) *user.User {
	return &user.User{
		ID:       id,
		Name:     "User Test",
		Email:    fmt.Sprintf("email%d@gmail.com", id),
		Password: "password",
		IsActive: false,
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

func TestStore(t *testing.T) {
	data := newData(1)
	db := getDB(t)
	defer clearAndClose(db, t)
	service := user.NewService(db, &user.Validator{})
	err := service.Store(data)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	service := user.NewService(db, &user.Validator{})
	data := newData(1)
	_ = service.Store(data)
	saved, err := service.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), saved.ID)
	assert.Equal(t, "email1@gmail.com", saved.Email)
	assert.Equal(t, false, saved.IsActive)
	assert.Equal(t, true, saved.IsAdmin)
}

func TestGetAll(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	service := user.NewService(db, &user.Validator{})
	b1 := newData(1)
	b2 := newData(2)
	_ = service.Store(b1)
	_ = service.Store(b2)
	saved, err := service.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(saved))
}

func TestUpdate(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	service := user.NewService(db, &user.Validator{})
	data := newData(1)
	_ = service.Store(data)
	t.Run("TestUpdate caminho feliz", func(t *testing.T) {
		saved, _ := service.Get(1)
		saved.Name = "UserTest2"
		saved.Email = "emailupdated@gmail.com"
		saved.IsActive = true
		saved.IsAdmin = false
		err := service.Update(saved)
		if err != nil {
			t.Fatalf("Erro atualizando %s", err.Error())
		}
		updated, _ := service.Get(1)
		assert.Equal(t, int64(1), updated.ID)
		assert.Equal(t, "emailupdated@gmail.com", updated.Email)
		assert.Equal(t, true, updated.IsActive)
		assert.Equal(t, false, updated.IsAdmin)

	})
	t.Run("TestUpdate erro de validação", func(t *testing.T) {
		e := newData(0)
		err := service.Update(e)
		if err == nil {
			t.Fatalf("Erro de validação")
		}
	})
}
