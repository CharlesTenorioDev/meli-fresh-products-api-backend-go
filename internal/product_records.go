package internal

import "time"

type ProductRecords struct {
	Id             int       `json:"id"`
	LastUpdateDate time.Time `json:"last_update_date"`
	PurchasePrice  float64   `json:"purchase_price"`
	SalePrice      float32   `json:"sale_price"`
	ProductID      int       `json:"product_id"`
}
type ProductRecordsJson struct {
	LastUpdateDate time.Time `json:"last_update_date"`
	PurchasePrice  float64   `json:"purchase_price"`
	SalePrice      float32   `json:"sale_price"`
	ProductID      int       `json:"product_id"`
}

type ProductRecordsJsonCount struct {
	ProductID    int    `json:"product_id"`
	Description  string `json:"description"`
	RecordsCount int    `json:"records_count"`
}

type ProductRecordsService interface {
	GetAll() ([]ProductRecords, error)
	GetByID(int) (ProductRecords, error)
	Create(ProductRecords) (ProductRecords, error)
}

type ProductRecordsRepository interface {
	FindAll() ([]ProductRecords, error)
	FindByID(int) (ProductRecords, error)
	Save(ProductRecords) (ProductRecords, error)
}
