package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	SectionNotFound            = errors.New("section not found")
	WarehouseNotFound          = errors.New("warehouse not found")
	ProductTypeNotFound        = errors.New("product_type not found")
	ProductNotFound            = errors.New("product not found")
	SectionAlreadyExists       = errors.New("section already exists")
	SectionNumberAlreadyInUse  = errors.New("section with given section number already registered")
	SectionUnprocessableEntity = errors.New("couldn't parse section")
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
		return internal.Section{}, SectionNotFound
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
		section.ProductTypeID <= 0 ||
		len(section.ProductBatches) == 0 {
		return errors.New("all fields must be valid and filled, unless otherwise stated")
	}

	return nil
}

func (s *SectionService) Save(section *internal.Section) error {
	if err := validateRequiredFields(*section); err != nil {
		return SectionUnprocessableEntity
	}

	err := s.rpS.SectionNumberExists(*section)
	if err != nil {
		return SectionNumberAlreadyInUse
	}

	_, err = s.rpW.FindByID(section.WarehouseID)
	if err != nil {
		return WarehouseNotFound
	}

	_, err = s.rpT.FindByID(section.ProductTypeID)
	if err != nil {
		return ProductTypeNotFound
	}

	for _, productBatch := range section.ProductBatches {
		_, err = s.rpP.FindByID(productBatch.ProductID)
		if err != nil {
			return ProductNotFound
		}
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
		return internal.Section{}, SectionNotFound
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
			return internal.Section{}, SectionNumberAlreadyInUse
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
			return internal.Section{}, WarehouseNotFound
		}
	}

	if _, ok := updates["product_type_id"]; ok {
		if err := processInt("product_type_id", &section.ProductTypeID); err != nil {
			return internal.Section{}, err
		}

		_, err = s.rpT.FindByID(section.ProductTypeID)
		if err != nil {
			return internal.Section{}, ProductTypeNotFound
		}
	}

	if productBatches, ok := updates["product_batches"]; ok {
		if batches, ok := productBatches.([]interface{}); ok {
			var batchs []internal.ProductBatch
			for _, batchItem := range batches {
				var batch internal.ProductBatch

				if batchMap, ok := batchItem.(map[string]interface{}); ok {
					if productId, ok := batchMap["product_id"].(float64); ok {
						batch.ProductID = int(productId)
					} else {
						return internal.Section{}, fmt.Errorf("product_id is required and must be a number")
					}

					_, err = s.rpP.FindByID(batch.ProductID)
					if err != nil {
						return internal.Section{}, ProductTypeNotFound
					}

					if quantity, ok := batchMap["quantity"].(float64); ok {
						batch.Quantity = int(quantity)
					} else {
						return internal.Section{}, fmt.Errorf("quantity is required and must be a number")
					}

					batchs = append(batchs, batch)
				} else {
					return internal.Section{}, fmt.Errorf("each item in product_batches must be an object")
				}
			}

			if len(batchs) > 0 {
				section.ProductBatches = batchs
			}

		} else {
			return internal.Section{}, fmt.Errorf("product_batches must be a list of objects")
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
		return SectionNotFound
	}

	err = s.rpS.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
