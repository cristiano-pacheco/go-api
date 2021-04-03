package validator_test

import (
	"testing"

	"github.com/cristiano-pacheco/go-api/core/validator"
	"github.com/stretchr/testify/assert"
)

func TestNotEmpty(t *testing.T) {
	assert.Equal(t, "cannot be empty", validator.NotEmpty("").Error())
	assert.Nil(t, validator.NotEmpty("teste"))
}

func TestMaxLength(t *testing.T) {
	assert.Equal(t, "is too long (maximum is 1 characters)", validator.MaxLength("bb", 1).Error())
	assert.Equal(t, "is too long (maximum is 3 characters)", validator.MaxLength("bbccc", 3).Error())
	assert.Nil(t, validator.MaxLength("bbccc", 6))
}

func TestMinLength(t *testing.T) {
	assert.Equal(t, "is too short (minimum is 3 characters)", validator.MinLength("bb", 3).Error())
	assert.Equal(t, "is too short (minimum is 6 characters)", validator.MinLength("bbccc", 6).Error())
	assert.Nil(t, validator.MinLength("bbccc", 5))
}

func TestIsEmail(t *testing.T) {
	assert.Equal(t, true, validator.IsEmail("teste@gmail.com"))
	assert.Equal(t, false, validator.IsEmail("teste"))
	assert.Equal(t, false, validator.IsEmail("teste@"))
}
