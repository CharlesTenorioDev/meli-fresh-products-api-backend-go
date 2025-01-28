package repository

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/loader"
)

func NewRepositorySection(dbPath string) (mp *SectionDB, err error) {
	bdSections := make(map[int]*internal.Section)

	sectionList, err := loader.ReadAllSectionsToFile(dbPath)
	if err != nil {
		return
	}

	for _, value := range sectionList {
		section := internal.Section{
			ID:                 value.ID,
			SectionNumber:      value.SectionNumber,
			CurrentTemperature: value.CurrentTemperature,
			MinimumTemperature: value.MinimumTemperature,
			CurrentCapacity:    value.CurrentCapacity,
			MinimumCapacity:    value.MinimumCapacity,
			MaximumCapacity:    value.MaximumCapacity,
			WarehouseID:        value.WarehouseID,
			ProductTypeID:      value.ProductTypeID,
		}

		bdSections[value.ID] = &section
	}

	mp = &SectionDB{
		DB:     bdSections,
		lastID: len(bdSections),
	}
	return
}

type SectionDB struct {
	DB     map[int]*internal.Section
	lastID int
}

func (r *SectionDB) FindAll() ([]internal.Section, error) {
	sections := make([]internal.Section, 0, len(r.DB))

	if len(r.DB) == 0 {
		return nil, errors.New("no sections available")
	}

	for _, section := range r.DB {
		sections = append(sections, *section)
	}

	return sections, nil
}

func (r *SectionDB) FindByID(id int) (internal.Section, error) {
	section, exists := r.DB[id]
	if !exists {
		return internal.Section{}, errors.New("section not found")
	}

	return *section, nil
}

func (r *SectionDB) ReportProducts() (int, error) {
	return 0, nil
}

func (r *SectionDB) ReportProductsByID(id int) (int, error) {
	return 0, nil
}

func (r *SectionDB) SectionNumberExists(section internal.Section) error {
	for _, value := range r.DB {
		if value.ID != section.ID && value.SectionNumber == section.SectionNumber {
			return errors.New("section with this section number already exists")
		}
	}

	return nil
}

func (r *SectionDB) Save(section *internal.Section) error {
	r.lastID++
	section.ID = r.lastID

	if _, exists := r.DB[section.ID]; exists {
		return errors.New("section already exists")
	}

	r.DB[section.ID] = section

	return nil
}

func (r *SectionDB) Update(section *internal.Section) error {
	if _, exists := r.DB[section.ID]; !exists {
		return errors.New("section not found")
	}

	r.DB[section.ID] = section

	return nil
}

func (r *SectionDB) Delete(id int) error {
	if _, exists := r.DB[id]; !exists {
		return errors.New("section not found")
	}

	delete(r.DB, id)

	return nil
}
