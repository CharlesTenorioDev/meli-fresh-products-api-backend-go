package internal

import "errors"

var (
	ErrProductBatchNotFound            = errors.New("product-batch not found")
	ErrProductBatchAlreadyExists       = errors.New("product-batch already exists")
	ErrProductBatchNumberAlreadyInUse  = errors.New("product-batch with given product-batch number already registered")
	ErrProductBatchUnprocessableEntity = errors.New("couldn't parse product-batch")
)

type ProductBatch struct {
	ID                 int     `json:"id"`
	BatchNumber        int     `json:"batch_number"`
	CurrentQuantity    int     `json:"current_quantity"`
	CurrentTemperature float64 `json:"current_temperature"`
	DueDate            string  `json:"due_date"`
	InitialQuantity    int     `json:"initial_quantity"`
	ManufacturingDate  string  `json:"manufacturing_date"`
	ManufacturingHour  int     `json:"manufacturing_hour"`
	MinumumTemperature float64 `json:"minumum_temperature"`
	ProductID          int     `json:"product_id"`
	SectionID          int     `json:"section_id"`
}

type ProductBatchRepository interface {
	FindByID(id int) (ProductBatch, error)
	Save(prodBatch *ProductBatch) error
	ProductBatchNumberExists(batchNumber int) (bool, error)
	ReportProducts() (prodBatches []ProductBatch, err error)
	ReportProductsByID(id int) (prodBatches []ProductBatch, err error)
}

type ProductBatchService interface {
	FindByID(id int) (ProductBatch, error)
	Save(prodBatch *ProductBatch) error
}

func (pb *ProductBatch) Ok() bool {
	if pb.BatchNumber <= 0 ||
		pb.CurrentQuantity < 0 ||
		pb.CurrentTemperature < -273 ||
		pb.DueDate == "" ||
		pb.InitialQuantity <= 0 ||
		pb.ManufacturingDate == "" ||
		pb.ManufacturingHour < 0 ||
		pb.MinumumTemperature < -273 ||
		pb.ProductID <= 0 ||
		pb.SectionID <= 0 {
		return false
	}

	return true
}
