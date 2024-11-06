package database

import (
	"github.com/antoniofmoliveira/apis/internal/entity"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		DB: db,
	}
}

func (r *ProductRepository) Create(product *entity.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepository) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	var products []entity.Product
	var err error
	if sort != "" && sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	if page != 0 && limit != 0 {
		err = r.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	} else {
		err = r.DB.Order("created_at " + sort).Find(&products).Error
	}
	return products, err
}

func (r *ProductRepository) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	err := r.DB.Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *ProductRepository) Update(product *entity.Product) (int64, error) {
	s := r.DB.Where("id = ?", product.ID).Updates(product)
	return s.RowsAffected, s.Error
}

func (r *ProductRepository) Delete(id string) (int64, error) {
	s := r.DB.Where("id = ?", id).Delete(&entity.Product{})
	return s.RowsAffected, s.Error
}
