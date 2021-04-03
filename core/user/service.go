package user

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // OK
	"golang.org/x/crypto/bcrypt"
)

// UseCase Define the interface with functions that will be used
type UseCase interface {
	GetAll() ([]*User, error)
	Get(ID int64) (*User, error)
	Store(u *User) error
	Update(u *User) error
	UpdatePassword(u *User) error
	Remove(ID int64) error
}

// Service define the struct for user service
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

// GetAll return all users from database
func (s *Service) GetAll() ([]*User, error) {
	var result []*User

	rows, err := s.DB.Query("select id, name, email, is_active, is_admin, created_at, updated_at from user")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsActive, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt)

		if err != nil {
			return nil, err
		}

		result = append(result, &u)
	}

	return result, nil
}

// Get the user from database
func (s *Service) Get(ID int64) (*User, error) {
	var u User

	stmt, err := s.DB.Prepare("select id, name, email, password, is_active, is_admin, created_at, updated_at from user where id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(ID).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsActive, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Store a user in the database
func (s *Service) Store(u *User) error {
	err := s.validator.validateUserCreationData(u)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into user(id, name, email, password, is_active, is_admin) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(u.ID, u.Name, u.Email, u.Password, u.IsActive, u.IsAdmin)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Update an user in the database
func (s *Service) Update(u *User) error {
	err := s.validator.validateUserUpdateData(u)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update user set name =?, email = ?, is_active = ?, is_admin = ? where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(u.Name, u.Email, u.IsActive, u.IsAdmin, u.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// UpdatePassword an user in the database
func (s *Service) UpdatePassword(u *User) error {
	err := s.validator.validateUserUpdatePasswordData(u)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("update user set password =? where id = ?")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(u.Password, u.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// Remove an user from the database
func (s *Service) Remove(ID int64) error {
	if ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("delete from user where id = ?", ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
