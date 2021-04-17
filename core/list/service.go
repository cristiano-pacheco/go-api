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
	Store(u *List) error
	Update(u *List) error
	UpdatePassword(u *List) error
	Remove(ID int64) error
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
	var u List

	stmt, err := s.DB.Prepare("select id, name, is_active, created_at, updated_at from list where id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&u.ID, &u.Name, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Store a record in the database
func (s *Service) Store(u *List) error {
	err := s.validator.validateCreationData(u)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into list(id, name, is_active, is_admin) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.ID, u.Name, u.IsActive)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Update an record in the database
func (s *Service) Update(u *List) error {
	err := s.validator.validateUpdateData(u)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update user set name =?, is_active = ?, where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(u.Name, u.IsActive, u.ID)
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
