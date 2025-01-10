package internal

import (
	"errors"
	"time"
)

var (
	ProductNotFound = errors.New("product not found")
)

type Product struct {
	Id                             int       `json:"id"`
	ProductCode                    string    `json:"product_code"`
	Description                    string    `json:"description"`
	Height                         float64   `json:"height"`
	Width                          float64   `json:"width"`
	NetWeight                      float64   `json:"net_weight"`
	ExpirationRate                 time.Time `json:"expiration_rate"`
	RecommendedFreezingTemperature float64   `json:"recommended_freezing_temperature"`
	FreezingRate                   float64   `json:"freezing_rate"`
	ProductTypeId                  int       `json:"product_type_id"`
	SellerId                       int       `json:"seller_id"`
}

type ProductService interface {
	GetAll() ([]Product, error)
	GetByID(id int) (Product, error)
	Create(Product) (Product, error)
	Update(Product) (Product, error)
	Delete(id int) error
}

type ProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id int) (Product, error)
	Save(Product) (Product, error)
	Update(Product) (Product, error)
	Delete(id int) error
}
