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

	if err := s.updateSectionNumber(updateSection.SectionNumber, &actualSection); err != nil {
		return internal.Section{}, err
	}

	if err := s.updateTemperature(&updateSection, &actualSection); err != nil {
		return internal.Section{}, err
	}

	if err := s.updateCapacity(&updateSection, &actualSection); err != nil {
		return internal.Section{}, err
	}

	if err := s.updateWarehouseAndProduct(&updateSection, &actualSection); err != nil {
		return internal.Section{}, err
	}

	if err := s.rpS.Update(&actualSection); err != nil {
		return internal.Section{}, err
	}

	return actualSection, nil
}

func (s *SectionService) updateSectionNumber(sectionNumber *int, actualSection *internal.Section) error {
	if sectionNumber != nil {
		countExists, err := s.rpS.SectionNumberExists(*sectionNumber)
		if err != nil {
			return err
		}

		if countExists {
			return internal.ErrSectionNumberAlreadyInUse
		}

		if *sectionNumber <= 0 {
			return internal.ErrSectionUnprocessableEntity
		}

		actualSection.SectionNumber = *sectionNumber
	}

	return nil
}

func (s *SectionService) updateTemperature(updateSection *internal.SectionPatch, actualSection *internal.Section) error {
	if updateSection.CurrentTemperature != nil {
		if *updateSection.CurrentTemperature < -273.15 {
			return internal.ErrSectionUnprocessableEntity
		}

		actualSection.CurrentTemperature = *updateSection.CurrentTemperature
	}

	if updateSection.MinimumTemperature != nil {
		if *updateSection.MinimumTemperature < -273.15 {
			return internal.ErrSectionUnprocessableEntity
		}

		actualSection.MinimumTemperature = *updateSection.MinimumTemperature
	}

	return nil
}

func (s *SectionService) updateCapacity(updateSection *internal.SectionPatch, actualSection *internal.Section) error {
	if updateSection.CurrentCapacity != nil {
		if *updateSection.CurrentCapacity < 0 {
			return internal.ErrSectionUnprocessableEntity
		}

		actualSection.CurrentCapacity = *updateSection.CurrentCapacity
	}

	if updateSection.MinimumCapacity != nil {
		if *updateSection.MinimumCapacity < 0 {
			return internal.ErrSectionUnprocessableEntity
		}

		actualSection.MinimumCapacity = *updateSection.MinimumCapacity
	}

	if updateSection.MaximumCapacity != nil {
		if *updateSection.MaximumCapacity < 0 {
			return internal.ErrSectionUnprocessableEntity
		}

		actualSection.MaximumCapacity = *updateSection.MaximumCapacity
	}

	return nil
}

func (s *SectionService) updateWarehouseAndProduct(updateSection *internal.SectionPatch, actualSection *internal.Section) error {
	if updateSection.WarehouseID != nil {
		_, err := s.rpW.FindByID(*updateSection.WarehouseID)
		if err != nil {
			return internal.ErrWarehouseRepositoryNotFound
		}

		actualSection.WarehouseID = *updateSection.WarehouseID
	}

	if updateSection.ProductTypeID != nil {
		_, err := s.rpW.FindByID(*updateSection.ProductTypeID)
		if err != nil {
			return internal.ErrProductTypeNotFound
		}

		actualSection.ProductTypeID = *updateSection.ProductTypeID
	}

	return nil
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
