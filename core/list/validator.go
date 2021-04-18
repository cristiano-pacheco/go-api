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

func (uv *Validator) validateListItemCreationData(li *ListItem) error {
	if li.ListID == 0 {
		return fmt.Errorf("invalid List ID")
	}

	if li.CategoryID == 0 {
		return fmt.Errorf("invalid Category ID")
	}

	err := validator.NotEmpty("name", li.Name)
	if err != nil {
		return err
	}

	return nil
}

func (uv *Validator) validateListItemUpdateData(li *ListItem) error {
	if li.ID == 0 {
		return fmt.Errorf("invalid List ID")
	}

	if li.CategoryID == 0 {
		return fmt.Errorf("invalid Category ID")
	}

	err := validator.NotEmpty("name", li.Name)
	if err != nil {
		return err
	}

	return nil
}
