package auth

import (
	"database/sql"

	"github.com/cristiano-pacheco/go-api/core/user"
)

// UseCase auth
type UseCase interface {
	IssueToken() (*Token, error)
}

// Service define the struct service
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

// IssueToken a new token
func (s *Service) IssueToken(email, password string) (*Token, error) {
	err := s.validator.validate(email, password)
	if err != nil {
		return nil, err
	}

	var u user.User
	var t Token

	stmt, err := s.DB.Prepare("select id, name, email, password, is_active, is_admin from user where email = ? and is_active = 1")
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRow(email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsActive, &u.IsAdmin)

	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	t.Token = u.Name

	return &t, nil
}
