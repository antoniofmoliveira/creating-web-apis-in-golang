package database

import (
	"testing"

	"github.com/antoniofmoliveira/apis/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)
	product, _ := entity.NewProduct("Product", 10.0)
	err = productRepository.Create(product)
	assert.Nil(t, err)

	product, err = productRepository.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, product.ID, product.ID)
	assert.Equal(t, "Product", product.Name)
	assert.Equal(t, 10.0, product.Price)
}

func TestFindAllProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)
	product, _ := entity.NewProduct("Product", 10.0)
	product2, _ := entity.NewProduct("Product 2", 20.0)
	err = productRepository.Create(product)
	assert.Nil(t, err)
	err = productRepository.Create(product2)
	assert.Nil(t, err)

	products, err := productRepository.FindAll(1, 2, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 2)
}

func TestFindAllProductsNoPagination(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)
	product, _ := entity.NewProduct("Product", 10.0)
	product2, _ := entity.NewProduct("Product 2", 20.0)
	err = productRepository.Create(product)
	assert.Nil(t, err)
	err = productRepository.Create(product2)
	assert.Nil(t, err)

	products, err := productRepository.FindAll(0, 0, "asc")
	assert.Nil(t, err)
	assert.Len(t, products, 2)
}
func TestUpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)
	product, _ := entity.NewProduct("Product", 10.0)
	err = productRepository.Create(product)
	assert.Nil(t, err)

	product, err = productRepository.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, product.ID, product.ID)
	assert.Equal(t, "Product", product.Name)
	assert.Equal(t, 10.0, product.Price)

	product.Name = "Product 2"
	rowsAffected, err := productRepository.Update(product)
	assert.Equal(t, int64(1), rowsAffected)
	assert.Nil(t, err)

	product, err = productRepository.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, product.ID, product.ID)
	assert.Equal(t, "Product 2", product.Name)
	assert.Equal(t, 10.0, product.Price)
}

func TestDeleteProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)
	product, _ := entity.NewProduct("Product", 10.0)
	err = productRepository.Create(product)
	assert.Nil(t, err)

	product, err = productRepository.FindByID(product.ID.String())
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, product.ID, product.ID)
	assert.Equal(t, "Product", product.Name)
	assert.Equal(t, 10.0, product.Price)

	rowsAffected, err := productRepository.Delete(product.ID.String())
	assert.Equal(t, int64(1), rowsAffected)
	assert.Nil(t, err)
}

func TestProductsNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)

	products, err := productRepository.FindAll(1, 2, "id")
	assert.Nil(t, err)
	assert.Len(t, products, 0)
}

func TestDeleteError(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Product{})

	productRepository := NewProductRepository(db)

	rowsAffected, err := productRepository.Delete("id")
	assert.Equal(t, int64(0), rowsAffected)
	assert.Nil(t, err)
}
