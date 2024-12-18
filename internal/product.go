package internal

import "time"

type Product struct {
	Id                             int     `json:"id"`
	ProductCode                    int     `json:"product_code"`
	Description                    string  `json:"description"`
	Height                         float64 `json:"height"`
	Width                          float64 `json:"width"`
	NetWeight                      float64 `json:"net_weight"`
	ExpirationRate                 time.Time `json:"expiration_rate"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	FreezingRate                   float64 `json:"freezing_rate"`
	ProductTypeId                  int     `json:"product_type_id"`
	SellerId                       int     `json:"seller_id"`
}


type ProductService interface {
	GetAll() (v map[int]Product, err error)
	

}

type ProductRepository interface {
	GetAll() (v map[int]Product, err error)
}