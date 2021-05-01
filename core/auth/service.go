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
	HasAccess(userID int, action string) (bool, error)
	IssueToken(email, password string) (*Token, error)
	GetUserPermissionsById(ID int64) (*UserPermission, error)
}

// Service define the struct service
type Service struct {
	DB        *sql.DB
	validator *Validator
	JWTHash   *jwt.HMACSHA
}

// NewService constructor
func NewService(db *sql.DB, v *Validator, jwtHash *jwt.HMACSHA) *Service {
	return &Service{
		DB:        db,
		validator: v,
		JWTHash:   jwtHash,
	}
}

// GetUserPermissionsById
func (s *Service) GetUserPermissionsById(ID int64) (*UserPermission, error) {
	stmt, err := s.DB.Prepare("select name from user where id = ?")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	au := &UserPermission{}
	au.ID = ID

	err = stmt.QueryRow(ID).Scan(&au.Name)

	if err != nil {
		return nil, err
	}

	var permissions []*Permission

	sql := `
		select p.name, p.action from user_permission up
		join permission p ON up.permission_id = p.id
		where user_id = ?
	`

	stmtPermissions, err := s.DB.Prepare(sql)

	if err != nil {
		return nil, err
	}

	defer stmtPermissions.Close()

	rows, err := stmtPermissions.Query(ID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p Permission
		err = rows.Scan(&p.Name, &p.Code)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, &p)
	}

	au.Permissions = permissions

	return au, nil
}

// HasAccess action
func (s *Service) HasAccess(userID int, action string) (bool, error) {
	stmt, err := s.DB.Prepare(
		"select count(1) from user_permission up join permission p on up.permission_id = p.id and p.action = ? where up.user_id = ?",
	)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var hasAccess = 0
	err = stmt.QueryRow(action, userID).Scan(&hasAccess)
	if err != nil {
		return false, err
	}
	if hasAccess == 1 {
		return true, nil
	}
	return false, nil
}

// IssueToken a new token
func (s *Service) IssueToken(email, password string) (*Token, error) {
	u, err := s.checkUserCredentials(email, password)
	if err != nil {
		return nil, err
	}

	var t Token
	now := time.Now()
	pl := CustomPayload{
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(time.Minute * 5)),
			IssuedAt:       jwt.NumericDate(now),
		},
		UserID: u.ID,
	}

	token, err := jwt.Sign(pl, s.JWTHash)
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
		return nil, err
	}

	err = stmt.QueryRow(email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsActive, &u.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("Invalid Credentials")
		}
		return nil, err
	}

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
