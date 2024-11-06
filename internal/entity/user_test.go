package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("John Doe", "j@j.com", "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "j@j.com", user.Email)
}

func TestNewUserWithEmptyEmail(t *testing.T) {
	_, err := NewUser("John Doe", "", "123456")
	assert.NotNil(t, err)
}

func TestNewUserWithInvalidEmail(t *testing.T) {
	_, err := NewUser("John Doe", "j@j", "123456")
	assert.NotNil(t, err)
}

func TestNewUserWithEmptyPassword(t *testing.T) {
	_, err := NewUser("John Doe", "j@j.com", "")
	assert.NotNil(t, err)
}

func TestNewUserWithEmptyName(t *testing.T) {
	_, err := NewUser("", "j@j.com", "123456")
	assert.NotNil(t, err)
}

func TestNewUserWithInvalidName(t *testing.T) {
	_, err := NewUser("", "j@j.com", "123456")
	assert.NotNil(t, err)
}

func TestNewUserWithInvalidPassword(t *testing.T) {
	password := strings.Repeat("a", 73)
	_, err := NewUser("John Doe", "j@j.com", password)
	assert.NotNil(t, err)
}

func TestValidatePassword(t *testing.T) {
	user, err := NewUser("John Doe", "j@j.com", "123456")
	assert.Nil(t, err)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("654321"))
	assert.NotEqual(t, "123456", user.Password)
}
