package list_test

import (
	"database/sql"
	"testing"

	"github.com/cristiano-pacheco/go-api/core/list"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func newData(id int64) *list.List {
	return &list.List{
		ID:       id,
		Name:     "List Test",
		IsActive: false,
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
	_, err = tx.Exec("delete from list")
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
	service := list.NewService(db, &list.Validator{})
	err := service.Store(data)
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	service := list.NewService(db, &list.Validator{})
	data := newData(1)
	_ = service.Store(data)
	saved, err := service.Get(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), saved.ID)
	assert.Equal(t, "List Test", saved.Name)
	assert.Equal(t, false, saved.IsActive)
}

func TestGetAll(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	service := list.NewService(db, &list.Validator{})
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
	service := list.NewService(db, &list.Validator{})
	data := newData(1)
	_ = service.Store(data)
	t.Run("TestUpdate caminho feliz", func(t *testing.T) {
		saved, _ := service.Get(1)
		saved.Name = "ListTest2"
		saved.IsActive = true
		err := service.Update(saved)
		if err != nil {
			t.Fatalf("Erro atualizando %s", err.Error())
		}
		updated, _ := service.Get(1)
		assert.Equal(t, int64(1), updated.ID)
		assert.Equal(t, "ListTest2", updated.Name)
		assert.Equal(t, true, updated.IsActive)

	})
	t.Run("TestUpdate erro de validação", func(t *testing.T) {
		e := newData(0)
		err := service.Update(e)
		if err == nil {
			t.Fatalf("Erro de validação")
		}
	})
}
