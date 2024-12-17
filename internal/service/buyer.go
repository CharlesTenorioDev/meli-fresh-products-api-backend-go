package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	BuyerNotFound = errors.New("buyer not found")
	BuyerAlreadyExists = errors.New("buyer already exists")
	BuyerUnprocessableEntity = errors.New("couldn't parse buyer")
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


func (s *BuyerServiceDefault) Save(buyer internal.Buyer) (err error) {
	all := s.repo.GetAll()
	if _, ok := all[buyer.ID]; ok {
		err = BuyerAlreadyExists
		return
	}

	ok := buyer.Parse()
	if !ok {
		err = BuyerUnprocessableEntity
		return
	}

	s.repo.Add(buyer)
	return
}


func (s* BuyerServiceDefault) Update(id int, buyerPatch internal.BuyerPatch) (err error) {
	all := s.repo.GetAll()
	_, ok := all[id]
	if !ok {
		err = BuyerNotFound
		return
	}

	s.repo.Update(id, buyerPatch)
	return
}

func (s* BuyerServiceDefault) Delete(id int) (err error) {
	all := s.repo.GetAll()
	_, ok := all[id]
	if !ok {
		err = BuyerNotFound
		return
	}

	s.repo.Delete(id)
	return
}
