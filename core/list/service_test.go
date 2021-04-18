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

func newItemData(id int64) *list.ListItem {
	return &list.ListItem{
		ID:         id,
		ListID:     1,
		CategoryID: 1,
		Name:       "List Item Test",
	}
}

func createCategoryAndList(db *sql.DB, t *testing.T) {
	tx, err := db.Begin()
	assert.Nil(t, err)
	_, err = tx.Exec("delete from list")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("delete from category")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("delete from list_item")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("insert into category (id, name, type) values (1, \"Teste\", 1)")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("insert into category (id, name, type) values (2, \"Teste 2\", 1)")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("insert into list (id, name, is_active) values (1, \"Teste\", 1)")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
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
	_, err = tx.Exec("delete from category")
	assert.Nil(t, err)
	if err != nil {
		tx.Rollback()
	}
	_, err = tx.Exec("delete from list_item")
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

func TestRemove(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	service := list.NewService(db, &list.Validator{})
	b1 := newData(1)
	b2 := newData(2)
	_ = service.Store(b1)
	_ = service.Store(b2)
	service.Remove(b1.ID)
	saved, err := service.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(saved))
}

func TestStoreItem(t *testing.T) {
	db := getDB(t)
	createCategoryAndList(db, t)
	data := newItemData(1)
	defer clearAndClose(db, t)
	service := list.NewService(db, &list.Validator{})
	err := service.StoreItem(data)
	assert.Nil(t, err)
}

func TestGetItem(t *testing.T) {
	db := getDB(t)
	createCategoryAndList(db, t)
	data := newItemData(1)
	defer clearAndClose(db, t)
	service := list.NewService(db, &list.Validator{})
	err := service.StoreItem(data)
	assert.Nil(t, err)
	saved, err := service.GetItem(1)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), saved.ID)
	assert.Equal(t, int64(1), saved.ListID)
	assert.Equal(t, int64(1), saved.CategoryID)
	assert.Equal(t, "List Item Test", saved.Name)
}

func TestGetAllItems(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	createCategoryAndList(db, t)
	service := list.NewService(db, &list.Validator{})
	b1 := newItemData(1)
	b2 := newItemData(2)
	_ = service.StoreItem(b1)
	_ = service.StoreItem(b2)
	saved, err := service.GetAllItems(1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(saved))
}

func TestUpdateItem(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	createCategoryAndList(db, t)
	service := list.NewService(db, &list.Validator{})
	data := newItemData(1)
	_ = service.StoreItem(data)
	t.Run("TestUpdateItem caminho feliz", func(t *testing.T) {
		saved, _ := service.GetItem(1)
		saved.Name = "ListItem2"
		saved.CategoryID = 2
		err := service.UpdateItem(saved)
		if err != nil {
			t.Fatalf("Erro atualizando %s", err.Error())
		}
		updated, _ := service.GetItem(1)
		assert.Equal(t, int64(1), updated.ID)
		assert.Equal(t, "ListItem2", updated.Name)
		assert.Equal(t, int64(2), updated.CategoryID)
		assert.Equal(t, int64(1), updated.ListID)
	})
	t.Run("TestUpdateItem erro de validação", func(t *testing.T) {
		e := newItemData(0)
		err := service.UpdateItem(e)
		if err == nil {
			t.Fatalf("Erro de validação")
		}
	})
}

func TestRemoveItem(t *testing.T) {
	db := getDB(t)
	defer clearAndClose(db, t)
	createCategoryAndList(db, t)
	service := list.NewService(db, &list.Validator{})
	b1 := newItemData(1)
	b2 := newItemData(2)
	_ = service.StoreItem(b1)
	_ = service.StoreItem(b2)
	service.RemoveItem(b1.ID)
	saved, err := service.GetAllItems(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(saved))
}
