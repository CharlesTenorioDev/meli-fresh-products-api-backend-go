package internal

type Section struct {
	ID                 int
	SectionNumber      int
	CurrentTemperature float64
	MinimumTemperature float64
	CurrentCapacity    int
	MinimumCapacity    int
	MaximumCapacity    int
	WarehouseID        int
	ProductTypeID      int
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
