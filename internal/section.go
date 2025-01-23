package internal

import "errors"

var (
	ErrSectionNotFound            = errors.New("section not found")
	ErrSectionAlreadyExists       = errors.New("section already exists")
	ErrSectionNumberAlreadyInUse  = errors.New("section with given section number already registered")
	ErrSectionUnprocessableEntity = errors.New("couldn't parse section")
)

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

type SectionPatch struct {
	SectionNumber      *int
	CurrentTemperature *float64
	MinimumTemperature *float64
	CurrentCapacity    *int
	MinimumCapacity    *int
	MaximumCapacity    *int
	WarehouseID        *int
	ProductTypeID      *int
}

type ReportProduct struct {
	SectionID     int `json:"section_id"`
	SectionNumber int `json:"section_number"`
	ProductsCount int `json:"products_count"`
}

type SectionRepository interface {
	FindAll() ([]Section, error)
	FindByID(id int) (Section, error)
	ReportProducts() ([]ReportProduct, error)
	ReportProductsByID(sectionID int) (ReportProduct, error)
	SectionNumberExists(sectionNumber int) (bool, error)
	Save(section *Section) error
	Update(section *Section) error
	Delete(id int) error
}

type SectionService interface {
	FindAll() ([]Section, error)
	FindByID(id int) (Section, error)
	ReportProducts() ([]ReportProduct, error)
	ReportProductsByID(sectionID int) (ReportProduct, error)
	Save(section *Section) error
	Update(id int, updateSection SectionPatch) (Section, error)
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
