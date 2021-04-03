package validator_test

import (
	"testing"

	"github.com/cristiano-pacheco/go-api/core/validator"
	"github.com/stretchr/testify/assert"
)

func TestNotEmpty(t *testing.T) {
	assert.Equal(t, "field cannot be empty", validator.NotEmpty("field", "").Error())
	assert.Nil(t, validator.NotEmpty("field", "teste"))
}

func TestMaxLength(t *testing.T) {
	assert.Equal(t, "field is too long (maximum is 1 characters)", validator.MaxLength("field", "bb", 1).Error())
	assert.Equal(t, "field is too long (maximum is 3 characters)", validator.MaxLength("field", "bbccc", 3).Error())
	assert.Nil(t, validator.MaxLength("field", "bbccc", 6))
}

func TestMinLength(t *testing.T) {
	assert.Equal(t, "field is too short (minimum is 3 characters)", validator.MinLength("field", "bb", 3).Error())
	assert.Equal(t, "field is too short (minimum is 6 characters)", validator.MinLength("field", "bbccc", 6).Error())
	assert.Nil(t, validator.MinLength("field", "bbccc", 5))
}

func TestEmail(t *testing.T) {
	assert.Equal(t, "field is not a valid email", validator.Email("field", "test").Error())
	assert.Equal(t, "field is not a valid email", validator.Email("field", "test@").Error())
	assert.Nil(t, validator.Email("field", "test@gmail.com"))
}

func TestIsEmail(t *testing.T) {
	assert.Equal(t, true, validator.IsEmail("teste@gmail.com"))
	assert.Equal(t, false, validator.IsEmail("teste"))
	assert.Equal(t, false, validator.IsEmail("teste@"))
}
