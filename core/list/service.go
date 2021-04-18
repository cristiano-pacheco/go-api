package list

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // OK
)

// UseCase Define the interface with functions that will be used
type UseCase interface {
	GetAll() ([]*List, error)
	Get(ID int64) (*List, error)
	Store(l *List) error
	Update(l *List) error
	Remove(ID int64) error
	GetAllItems(listID int64) ([]*ListItem, error)
	GetItem(ID int64) (*ListItem, error)
	StoreItem(li *ListItem) error
	UpdateItem(li *ListItem) error
	RemoveItem(ID int64) error
}

// Service define the struct for service
type Service struct {
	DB        *sql.DB
	validator *Validator
}

// NewService constructor
func NewService(db *sql.DB, v *Validator) *Service {
	return &Service{
		DB:        db,
		validator: v,
	}
}

// GetAll return all records from the database
func (s *Service) GetAll() ([]*List, error) {
	var result []*List

	rows, err := s.DB.Query("select id, name, is_active, created_at, updated_at from list")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u List
		err := rows.Scan(&u.ID, &u.Name, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

		if err != nil {
			return nil, err
		}

		result = append(result, &u)
	}

	return result, nil
}

// Get the records from the database
func (s *Service) Get(ID int64) (*List, error) {
	var l List

	stmt, err := s.DB.Prepare("select id, name, is_active, created_at, updated_at from list where id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&l.ID, &l.Name, &l.IsActive, &l.CreatedAt, &l.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &l, nil
}

// Store a record in the database
func (s *Service) Store(l *List) error {
	err := s.validator.validateCreationData(l)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into list(id, name, is_active) values (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(l.ID, l.Name, l.IsActive)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Update an record in the database
func (s *Service) Update(l *List) error {
	err := s.validator.validateUpdateData(l)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update list set name =?, is_active = ? where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(l.Name, l.IsActive, l.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Remove an record from the database
func (s *Service) Remove(ID int64) error {
	if ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from list where id = ?", ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// GetAllItems return all records from the database
func (s *Service) GetAllItems(listID int64) ([]*ListItem, error) {
	var result []*ListItem

	sql := `
		select li.id, li.list_id, li.category_id, c.name as category_name, li.name, li.created_at, li.updated_at 
		from list_item as li
		left join category c on li.category_id = c.id
		where li.list_id = ?
	`

	stmt, err := s.DB.Prepare(sql)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(listID)

	for rows.Next() {
		var li ListItem
		err := rows.Scan(&li.ID, &li.ListID, &li.CategoryID, &li.CategoryName, &li.Name, &li.CreatedAt, &li.UpdatedAt)

		if err != nil {
			return nil, err
		}

		result = append(result, &li)
	}

	return result, nil
}

// GetItem the record from the database
func (s *Service) GetItem(ID int64) (*ListItem, error) {
	var li ListItem

	sql := `
		select li.id, li.list_id, li.category_id, c.name as category_name, li.name, li.created_at, li.updated_at 
		from list_item as li
		left join category c on li.category_id = c.id
		where li.id = ?
	`

	stmt, err := s.DB.Prepare(sql)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&li.ID, &li.ListID, &li.CategoryID, &li.CategoryName, &li.Name, &li.CreatedAt, &li.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &li, nil
}

// Store a record in the database
func (s *Service) StoreItem(li *ListItem) error {
	err := s.validator.validateListItemCreationData(li)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into list_item (id, list_id, category_id, name) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(li.ID, li.ListID, li.CategoryID, li.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Update an record in the database
func (s *Service) UpdateItem(li *ListItem) error {
	err := s.validator.validateListItemUpdateData(li)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update list_item set category_id=?, name =? where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(li.CategoryID, li.Name, li.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Remove an record from the database
func (s *Service) RemoveItem(ID int64) error {
	if ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from list_item where id = ?", ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
