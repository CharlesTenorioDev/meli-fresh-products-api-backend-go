package service

import (
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

	countExists, err := s.rpS.SectionNumberExists(section.SectionNumber)
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

func (s *SectionService) Update(id int, updateSection internal.SectionPatch) (internal.Section, error) {
	actualSection, err := s.FindByID(id)
	if err != nil {
		return internal.Section{}, err
	}

	if updateSection.SectionNumber != nil {
		countExists, err := s.rpS.SectionNumberExists(*updateSection.SectionNumber)
		if err != nil || countExists {
			return internal.Section{}, internal.ErrSectionNumberAlreadyInUse
		}

		actualSection.SectionNumber = *updateSection.SectionNumber
		if actualSection.SectionNumber <= 0 {
			return internal.Section{}, internal.ErrSectionUnprocessableEntity
		}
	}

	if updateSection.CurrentTemperature != nil {
		actualSection.CurrentTemperature = *updateSection.CurrentTemperature
		if actualSection.CurrentTemperature < -273.15 {
			return internal.Section{}, internal.ErrSectionUnprocessableEntity
		}
	}

	if updateSection.MinimumTemperature != nil {
		actualSection.MinimumTemperature = *updateSection.MinimumTemperature
		if actualSection.MinimumTemperature < -273.15 {
			return internal.Section{}, internal.ErrSectionUnprocessableEntity
		}
	}

	if updateSection.CurrentCapacity != nil {
		actualSection.CurrentCapacity = *updateSection.CurrentCapacity
		if actualSection.CurrentCapacity < 0 {
			return internal.Section{}, internal.ErrSectionUnprocessableEntity
		}
	}

	if updateSection.MinimumCapacity != nil {
		actualSection.MinimumCapacity = *updateSection.MinimumCapacity
		if actualSection.MinimumCapacity < 0 {
			return internal.Section{}, internal.ErrSectionUnprocessableEntity
		}
	}

	if updateSection.MaximumCapacity != nil {
		actualSection.MaximumCapacity = *updateSection.MaximumCapacity
		if actualSection.MaximumCapacity < 0 {
			return internal.Section{}, internal.ErrSectionUnprocessableEntity
		}
	}

	if updateSection.WarehouseID != nil {
		_, err := s.rpW.FindByID(*updateSection.WarehouseID)
		if err != nil {
			return internal.Section{}, internal.ErrWarehouseRepositoryNotFound
		}

		actualSection.WarehouseID = *updateSection.WarehouseID
	}

	if updateSection.ProductTypeID != nil {
		_, err := s.rpW.FindByID(*updateSection.ProductTypeID)
		if err != nil {
			return internal.Section{}, internal.ErrProductTypeNotFound
		}

		actualSection.ProductTypeID = *updateSection.ProductTypeID
	}

	err = s.rpS.Update(&actualSection)
	if err != nil {
		return internal.Section{}, err
	}

	return actualSection, nil
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
