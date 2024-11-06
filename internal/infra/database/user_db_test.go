package database

import (
	"testing"

	"github.com/antoniofmoliveira/apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.User{})

	userRepository := NewUserRepository(db)
	user, _ := entity.NewUser("John Doe", "j@j.com", "123456")
	err = userRepository.Create(user)
	assert.Nil(t, err)

	user, err = userRepository.FindByEmail("j@j.com")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Email, "j@j.com")
	assert.Equal(t, user.Name, "John Doe")
	assert.NotNil(t, user.Password)
}

func TestFindByEmail(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.User{})

	userRepository := NewUserRepository(db)
	user, _ := entity.NewUser("John Doe", "j@j.com", "123456")
	err = userRepository.Create(user)
	assert.Nil(t, err)

	user, err = userRepository.FindByEmail("j@j.com")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.Email, "j@j.com")
	assert.Equal(t, user.Name, "John Doe")
	assert.NotNil(t, user.Password)
}

func TestFindByEmailNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.User{})

	userRepository := NewUserRepository(db)
	_, err = userRepository.FindByEmail("j@j.com")
	assert.NotNil(t, err)
}
