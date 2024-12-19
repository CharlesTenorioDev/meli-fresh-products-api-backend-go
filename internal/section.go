package internal

type ProductBatch struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Section struct {
	ID                 int            `json:"id"`
	SectionNumber      int            `json:"section_number"`
	CurrentTemperature float64        `json:"current_temperature"`
	MinimumTemperature float64        `json:"minimum_temperature"`
	CurrentCapacity    int            `json:"current_capacity"`
	MinimumCapacity    int            `json:"minimum_capacity"`
	MaximumCapacity    int            `json:"maximum_capacity"`
	WarehouseID        int            `json:"warehouse_id"`
	ProductTypeID      int            `json:"product_type_id"`
	ProductBatches     []ProductBatch `json:"product_batches"`
}

type SectionRepository interface {
	FindAll() ([]Section, error)
	FindByID(id int) (Section, error)
	SectionNumberExists(section Section) error
	Save(section *Section) error
	Update(section *Section) error
	Delete(id int) error
}

type SectionService interface {
	FindAll() ([]Section, error)
	FindByID(id int) (Section, error)
	Save(section *Section) error
	Update(id int, updates map[string]interface{}) (Section, error)
	Delete(id int) error
}
