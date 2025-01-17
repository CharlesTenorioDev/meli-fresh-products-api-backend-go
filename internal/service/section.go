package service

import (
	"fmt"
	"strconv"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewServiceSection(rpSection internal.SectionRepository, rpProductType internal.ProductTypeRepository, rpProduct internal.ProductRepository, rpWareHouse internal.WarehouseRepository) *SectionService {
	return &SectionService{
		rpS: rpSection,
		rpP: rpProduct,
		rpT: rpProductType,
		rpW: rpWareHouse,
	}
}

type SectionService struct {
	rpS internal.SectionRepository
	rpP internal.ProductRepository
	rpT internal.ProductTypeRepository
	rpW internal.WarehouseRepository
}

func (s *SectionService) FindAll() ([]internal.Section, error) {
	sections, err := s.rpS.FindAll()
	if err != nil {
		return nil, err
	}

	return sections, nil
}

func (s *SectionService) FindByID(id int) (internal.Section, error) {
	section, err := s.rpS.FindByID(id)
	if err != nil {
		return internal.Section{}, internal.ErrSectionNotFound
	}

	return section, nil
}

func (s *SectionService) ReportProducts() ([]internal.ReportProduct, error) {
	reportProducts, err := s.rpS.ReportProducts()
	if err != nil {
		return nil, err
	}

	return reportProducts, nil
}

func (s *SectionService) ReportProductsByID(sectionID int) (internal.ReportProduct, error) {
	_, err := s.rpS.FindByID(sectionID)
	if err != nil {
		return internal.ReportProduct{}, internal.ErrSectionNotFound
	}

	reportProduct, err := s.rpS.ReportProductsByID(sectionID)
	if err != nil {
		return internal.ReportProduct{}, err
	}

	return reportProduct, nil
}

func (s *SectionService) Save(section *internal.Section) error {
	if ok := section.Ok(); !ok {
		return internal.ErrSectionUnprocessableEntity
	}

	countExists, err := s.rpS.SectionNumberExists(*section)
	if err != nil || countExists {
		return internal.ErrSectionNumberAlreadyInUse
	}

	_, err = s.rpW.FindByID(section.WarehouseID)
	if err != nil {
		return internal.ErrWarehouseRepositoryNotFound
	}

	_, err = s.rpT.FindByID(section.ProductTypeID)
	if err != nil {
		return internal.ErrProductTypeNotFound
	}

	err = s.rpS.Save(section)
	if err != nil {
		return err
	}

	return nil
}

func (s *SectionService) Update(id int, updates map[string]interface{}) (internal.Section, error) {
	section, err := s.FindByID(id)
	if err != nil {
		return internal.Section{}, internal.ErrSectionNotFound
	}

	processInt := func(key string, target *int) error {
		if val, ok := updates[key]; ok {
			switch v := val.(type) {
			case string:
				value, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("invalid value for %s: %v", key, err)
				}
				*target = value
			case float64:
				*target = int(v)
			default:
				return fmt.Errorf("invalid type for %s: expected string or float64, got %T", key, v)
			}
		}
		return nil
	}

	processFloat := func(key string, target *float64) error {
		if val, ok := updates[key]; ok {
			switch v := val.(type) {
			case string:
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf("invalid value for %s: %v", key, err)
				}
				*target = value
			case float64:
				*target = v
			default:
				return fmt.Errorf("invalid type for %s: expected string or float64, got %T", key, v)
			}
		}
		return nil
	}

	if _, ok := updates["section_number"]; ok {
		if err := processInt("section_number", &section.SectionNumber); err != nil {
			return internal.Section{}, err
		}

		countExists, err := s.rpS.SectionNumberExists(section)
		if err != nil || countExists {
			return internal.Section{}, internal.ErrSectionNumberAlreadyInUse
		}
	}

	if err := processFloat("current_temperature", &section.CurrentTemperature); err != nil {
		return internal.Section{}, err
	}

	if err := processFloat("minimum_temperature", &section.MinimumTemperature); err != nil {
		return internal.Section{}, err
	}

	if err := processInt("current_capacity", &section.CurrentCapacity); err != nil {
		return internal.Section{}, err
	}

	if err := processInt("minimum_capacity", &section.MinimumCapacity); err != nil {
		return internal.Section{}, err
	}

	if err := processInt("maximum_capacity", &section.MaximumCapacity); err != nil {
		return internal.Section{}, err
	}

	if _, ok := updates["warehouse_id"]; ok {
		if err := processInt("warehouse_id", &section.WarehouseID); err != nil {
			return internal.Section{}, err
		}

		_, err = s.rpW.FindByID(section.WarehouseID)
		if err != nil {
			return internal.Section{}, internal.ErrWarehouseRepositoryNotFound
		}
	}

	if _, ok := updates["product_type_id"]; ok {
		if err := processInt("product_type_id", &section.ProductTypeID); err != nil {
			return internal.Section{}, err
		}

		_, err = s.rpT.FindByID(section.ProductTypeID)
		if err != nil {
			return internal.Section{}, internal.ErrProductTypeNotFound
		}
	}

	err = s.rpS.Update(&section)
	if err != nil {
		return internal.Section{}, err
	}

	return section, nil
}

func (s *SectionService) Delete(id int) error {
	_, err := s.FindByID(id)
	if err != nil {
		return internal.ErrSectionNotFound
	}

	err = s.rpS.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
