package entity

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/antoniofmoliveira/apis/pkg/entity"
)

type User struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
}

var (
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	emailRegex         = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func NewUser(name, email, password string) (*User, error) {
	if password == "" {
		return nil, ErrInvalidPassword
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if email == "" || !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       entity.NewId(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
