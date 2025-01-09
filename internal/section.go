package internal

type Section struct {
	ID                 int     `json:"id"`
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}

type SectionRepository interface {
	FindAll() ([]Section, error)
	FindByID(id int) (Section, error)
	ReportProducts() (prodBatchs []ProductBatch, err error)
	ReportProductsByID(id int) (prodBatchs []ProductBatch, err error)
	SectionNumberExists(section Section) error
	Save(section *Section) error
	Update(section *Section) error
	Delete(id int) error
}

type SectionService interface {
	FindAll() ([]Section, error)
	FindByID(id int) (Section, error)
	ReportProducts() (prodBatchs []ProductBatch, err error)
	ReportProductsByID(id int) (prodBatchs []ProductBatch, err error)
	Save(section *Section) error
	Update(id int, updates map[string]interface{}) (Section, error)
	Delete(id int) error
}

func (s *Section) Ok() bool {
	if s.SectionNumber <= 0 ||
		s.CurrentTemperature < -273.15 ||
		s.MinimumTemperature < -273.15 ||
		s.CurrentCapacity < 0 ||
		s.MinimumCapacity < 0 ||
		s.MaximumCapacity < 0 ||
		s.WarehouseID <= 0 ||
		s.ProductTypeID <= 0 {
		return false
	}
	return true
}
