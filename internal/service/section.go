package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	errorss "github.com/meli-fresh-products-api-backend-t1/internal/errors"
)

func NewServiceSection(rpSection internal.SectionRepository, rpWareHouse internal.WarehouseRepository) *SectionService {
	return &SectionService{
		rpS: rpSection,
		rpW: rpWareHouse,
	}
}

type SectionService struct {
	rpS internal.SectionRepository
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
		return internal.Section{}, err
	}

	return section, nil
}

func validateRequiredFields(section internal.Section) error {
	if section.SectionNumber <= 0 ||
		section.CurrentTemperature < -273.15 ||
		section.MinimumTemperature < -273.15 ||
		section.CurrentCapacity < 0 ||
		section.MinimumCapacity < 0 ||
		section.MaximumCapacity < 0 ||
		section.WarehouseID <= 0 ||
		section.ProductTypeID <= 0 {
		return errorss.NewUnprocessableEntity("all fields must be valid and filled, unless otherwise stated")
	}

	return nil
}

func (s *SectionService) Save(section *internal.Section) error {
	if err := validateRequiredFields(*section); err != nil {
		return err
	}

	err := s.rpS.SectionNumberExists(*section)
	if err != nil {
		return errorss.NewConflictError(err.Error())
	}

	_, err = s.rpW.FindByID(section.WarehouseID)
	if err != nil {
		return errorss.NewNotFound("warehouse not found")
	}

	err = s.rpS.Save(section)
	if err != nil {
		return err
	}

	return err
}

func (s *SectionService) Update(id int, updates map[string]interface{}) (internal.Section, error) {
	section, err := s.FindByID(id)
	if err != nil {
		return internal.Section{}, errors.New("section not found")
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

		err := s.rpS.SectionNumberExists(section)
		if err != nil {
			return internal.Section{}, errorss.NewConflictError(err.Error())
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

	if _, ok := updates["section_number"]; ok {
		if err := processInt("warehouse_id", &section.WarehouseID); err != nil {
			return internal.Section{}, err
		}

		_, err = s.rpW.FindByID(section.WarehouseID)
		if err != nil {
			return internal.Section{}, errorss.NewNotFound("warehouse not found")
		}
	}

	if err := processInt("product_type_id", &section.ProductTypeID); err != nil {
		return internal.Section{}, err
	}

	err = s.rpS.Update(&section)
	if err != nil {
		return internal.Section{}, err
	}

	return section, nil
}

func (s *SectionService) Delete(id int) (err error) {
	err = s.rpS.Delete(id)
	if err != nil {
		return err
	}

	return
}
