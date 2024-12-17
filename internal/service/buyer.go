package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

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
