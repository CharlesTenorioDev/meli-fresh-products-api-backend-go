package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	BuyerNotFound = errors.New("buyer not found")
)

type BuyerServiceDefault struct {
	repo internal.BuyerRepository
}

func NewBuyerService(r internal.BuyerRepository) *BuyerServiceDefault {
	return &BuyerServiceDefault{
		repo: r,
	}
}

func (s *BuyerServiceDefault) GetAll() map[int]internal.Buyer {
	all := s.repo.GetAll()

	return all
}

func (s *BuyerServiceDefault) FindByID(id int) (b internal.Buyer, err error) {
	all := s.repo.GetAll()
	b, ok := all[id]
	if !ok {
		err = BuyerNotFound
	}

	return
}
