package service

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewServiceProductBatch(rpProductBatch internal.ProductBatchRepository, rpSection internal.SectionRepository, rpProduct internal.ProductRepository) *ProductBatchService {
	return &ProductBatchService{
		rpB: rpProductBatch,
		rpS: rpSection,
		rpP: rpProduct,
	}
}

type ProductBatchService struct {
	rpB internal.ProductBatchRepository
	rpS internal.SectionRepository
	rpP internal.ProductRepository
}

func (s *ProductBatchService) FindByID(id int) (internal.ProductBatch, error) {
	prodBatch, err := s.rpB.FindByID(id)
	if err != nil {
		return internal.ProductBatch{}, internal.ProductBatchNotFound
	}

	return prodBatch, nil
}

func (s *ProductBatchService) Save(prodBatch *internal.ProductBatch) error {
	if ok := prodBatch.Ok(); !ok {
		return internal.ProductBatchUnprocessableEntity
	}

	countExists, err := s.rpB.ProductBatchNumberExists(prodBatch.BatchNumber)
	if err != nil || countExists {
		return internal.ProductBatchNumberAlreadyInUse
	}

	_, err = s.rpP.FindByID(prodBatch.ProductId)
	if err != nil {
		return internal.ErrProductNotFound
	}

	_, err = s.rpS.FindByID(prodBatch.SectionId)
	if err != nil {
		return internal.SectionNotFound
	}

	err = s.rpB.Save(prodBatch)
	if err != nil {
		return err
	}

	return nil
}
