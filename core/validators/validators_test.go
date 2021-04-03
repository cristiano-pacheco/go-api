package validators_test

import (
	"testing"

	"github.com/cristiano-pacheco/go-api/core/validators"
	"github.com/stretchr/testify/assert"
)

func TestNotEmpty(t *testing.T) {
	assert.Equal(t, "cannot be empty", validators.NotEmpty("").Error())
	assert.Nil(t, validators.NotEmpty("teste"))
}

func TestMaxLength(t *testing.T) {
	assert.Equal(t, "is too long (maximum is 1 characters)", validators.MaxLength("bb", 1).Error())
	assert.Equal(t, "is too long (maximum is 3 characters)", validators.MaxLength("bbccc", 3).Error())
	assert.Nil(t, validators.MaxLength("bbccc", 6))
}

func TestMinLength(t *testing.T) {
	assert.Equal(t, "is too short (minimum is 3 characters)", validators.MinLength("bb", 3).Error())
	assert.Equal(t, "is too short (minimum is 6 characters)", validators.MinLength("bbccc", 6).Error())
	assert.Nil(t, validators.MinLength("bbccc", 5))
}

func TestIsEmail(t *testing.T) {
	assert.Equal(t, true, validators.IsEmail("teste@gmail.com"))
	assert.Equal(t, false, validators.IsEmail("teste"))
	assert.Equal(t, false, validators.IsEmail("teste@"))
}
