package entity

import (
	"errors"
	"time"

	"github.com/antoniofmoliveira/apis/pkg/entity"
)

var (
	ErrIDIsRequired    = errors.New("id is required")
	ErrInvalidID       = errors.New("invalid id")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceIsRequired = errors.New("price is required")
	ErrInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	ID        entity.ID `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

func NewProduct(name string, price float64) (*Product, error) {
	p := &Product{
		ID:        entity.NewId(),
		Name:      name,
		Price:     price,
		CreatedAt: time.Now(),
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseId(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceIsRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
