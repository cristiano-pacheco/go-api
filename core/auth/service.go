package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cristiano-pacheco/go-api/core/user"
	"github.com/gbrlsnchs/jwt/v3"
	"golang.org/x/crypto/bcrypt"
)

// UseCase auth
type UseCase interface {
	IssueToken() (*Token, error)
}

// Service define the struct service
type Service struct {
	DB        *sql.DB
	validator *Validator
	jwtKey    string
}

// NewService constructor
func NewService(db *sql.DB, v *Validator, jk string) *Service {
	return &Service{
		DB:        db,
		validator: v,
		jwtKey:    jk,
	}
}

// IssueToken a new token
func (s *Service) IssueToken(email, password string) (*Token, error) {
	u, err := s.checkUserCredentials(email, password)
	if err != nil {
		return nil, err
	}

	var t Token
	type CustomPayload struct {
		jwt.Payload
		UserID int64 `json:"user_id"`
	}

	var hs = jwt.NewHS256([]byte("secret"))
	now := time.Now()
	pl := CustomPayload{
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(time.Minute * 5)),
			IssuedAt:       jwt.NumericDate(now),
		},
		UserID: u.ID,
	}

	token, err := jwt.Sign(pl, hs)
	if err != nil {
		return nil, err
	}

	t.Token = string(token)

	return &t, nil
}

func (s *Service) checkUserCredentials(email, password string) (*user.User, error) {
	err := s.validator.validate(email, password)
	if err != nil {
		return nil, err
	}

	var u user.User

	stmt, err := s.DB.Prepare("select id, name, email, password, is_active, is_admin from user where email = ? and is_active = 1")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Invalid Credentials")
		}
		return nil, err
	}

	err = stmt.QueryRow(email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsActive, &u.IsAdmin)

	defer stmt.Close()

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, fmt.Errorf("Invalid Credentials")
		}
		return nil, err
	}

	return &u, nil
}
