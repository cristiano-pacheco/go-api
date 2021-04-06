package authorization

import "database/sql"

// UseCase auth
type UseCase interface {
	HasAccess(userID int64, action string) (bool, error)
}

// Service define the struct service
type Service struct {
	DB *sql.DB
}

// NewService constructor
func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

// HasAccess action
func (s *Service) HasAccess(userID int64, action string) (bool, error) {
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
