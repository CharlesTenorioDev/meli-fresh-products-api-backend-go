package internal

import (
	"errors"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductConflit       = errors.New("product conflict")
	ErrProductConflitEntity = errors.New("product has a dependency on another entity and needs to be deleted first")

	ErrProductTypeIDNotFound      = errors.New("product-type ID not found")
	ErrSellerIdNotFound           = errors.New("seller ID not found")
	ErrProductCodeAlreadyExists   = errors.New("product-type already exists")
	ErrProductUnprocessableEntity = errors.New("all fields must be valid and filled")
	ErrProductBadRequest          = errors.New("invalid syntax")
)

type Product struct {
	ID                             int     `json:"id"`
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	Width                          float64 `json:"width"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductTypeID                  int     `json:"product_type_id"`
	SellerID                       int     `json:"seller_id"`
}

type ProductJSONPost struct {
	ProductCode                    string  `json:"product_code"`
	Description                    string  `json:"description"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 float64 `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	Width                          float64 `json:"width"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductTypeID                  int     `json:"product_type_id"`
	SellerID                       int     `json:"seller_id"`
}

type ProductService interface {
	GetAll() ([]Product, error)
	GetByID(id int) (Product, error)
	Create(Product) (Product, error)
	Update(Product) (Product, error)
	Delete(id int) error
	GetByIDRecord(int) (ProductRecordsJSONCount, error)
	GetAllRecord() ([]ProductRecordsJSONCount, error)
}

type ProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id int) (Product, error)
	Save(Product) (Product, error)
	Update(Product) (Product, error)
	Delete(id int) error
	FindByIDRecord(int) (ProductRecordsJSONCount, error)
	FindAllRecord() ([]ProductRecordsJSONCount, error)
}
