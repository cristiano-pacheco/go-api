package auth

import (
	"github.com/cristiano-pacheco/go-api/core/validator"
)

// Validator struct
type Validator struct{}

func (uv *Validator) validate(email, password string) error {
	err := validator.NotEmpty("email", email)
	if err != nil {
		return err
	}

	err = validator.Email("email", email)
	if err != nil {
		return err
	}

	err = validator.NotEmpty("password", password)
	if err != nil {
		return err
	}

	return nil
}
