package internal

type ProductBatch struct {
	ID                 int
	BatchNumber        int
	CurrentQuantity    int
	CurrentTemperature float64
	DueDate            string
	InitialQuantity    int
	ManufacturingDate  string
	ManufacturingHour  int
	MinumumTemperature float64
	ProductId          int
	SectionId          int
}

type ProductBatchRepository interface {
	FindByID(id int) (ProductBatch, error)
	Save(prodBatch *ProductBatch) error
	ProductBatchNumberExists(batchNumber int) (bool, error)
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
		pb.ProductId <= 0 ||
		pb.SectionId <= 0 {
		return false
	}
	return true
}
