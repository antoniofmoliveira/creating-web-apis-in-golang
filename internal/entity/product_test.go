package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Product", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "Product", product.Name)
	assert.Equal(t, 10.0, product.Price)
}
func TestNewProductWithEmptyName(t *testing.T) {
	p, err := NewProduct("", 10.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrNameIsRequired, err)
}

func TestNewProductWithInvalidPrice(t *testing.T) {
	p, err := NewProduct("Product", -10.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrInvalidPrice, err)
}

func TestNewProductWith0Price(t *testing.T) {
	p, err := NewProduct("Product", 0.0)
	assert.Nil(t, p)
	assert.Equal(t, ErrPriceIsRequired, err)
}

func TestValidate(t *testing.T) {
	p, err := NewProduct("Product", 10.0)
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Nil(t, p.Validate())
}
