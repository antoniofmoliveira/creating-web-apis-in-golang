package database

import "github.com/antoniofmoliveira/apis/internal/entity"

type UserRepositoryInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}

type ProductRepositoryInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Update(product *entity.Product) (int64, error)
	Delete(id string) (int64, error)
}
