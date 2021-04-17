package list

import (
	"fmt"

	"github.com/cristiano-pacheco/go-api/core/validator"
)

// Validator struct
type Validator struct{}

func (uv *Validator) validateCreationData(l *List) error {
	err := validator.NotEmpty("name", l.Name)
	if err != nil {
		return err
	}

	return nil
}

func (uv *Validator) validateUpdateData(l *List) error {
	if l.ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	err := validator.NotEmpty("name", l.Name)
	if err != nil {
		return err
	}

	return nil
}
