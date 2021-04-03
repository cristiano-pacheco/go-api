package user

import (
	"fmt"

	"github.com/cristiano-pacheco/go-api/core/validator"
)

// Validator struct
type Validator struct{}

func (uv *Validator) validateUserCreationData(u *User) error {
	err := validator.NotEmpty("name", u.Name)
	if err != nil {
		return err
	}

	err = validator.NotEmpty("email", u.Email)
	if err != nil {
		return err
	}

	err = validator.Email("email", u.Email)
	if err != nil {
		return err
	}

	err = validator.NotEmpty("password", u.Password)
	if err != nil {
		return err
	}

	return nil
}

func (uv *Validator) validateUserUpdateData(u *User) error {
	if u.ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	err := validator.NotEmpty("name", u.Name)
	if err != nil {
		return err
	}

	err = validator.NotEmpty("email", u.Email)
	if err != nil {
		return err
	}

	err = validator.Email("email", u.Email)
	if err != nil {
		return err
	}

	return nil
}

func (uv *Validator) validateUserUpdatePasswordData(u *User) error {
	if u.ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	err := validator.NotEmpty("password", u.Name)
	if err != nil {
		return err
	}

	return nil
}
