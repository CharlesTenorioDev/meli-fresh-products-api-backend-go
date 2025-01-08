package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	BuyerNotFound            = errors.New("buyer not found")
	BuyerAlreadyExists       = errors.New("buyer already exists")
	CardNumberAlreadyInUse   = errors.New("buyer with given card number already registered")
	BuyerUnprocessableEntity = errors.New("couldn't parse buyer")
)

func cardNumberIdAlreadyInUse(cardNumber string, buyers map[int]internal.Buyer) bool {
	for _, b := range buyers {
		if b.CardNumberId == cardNumber {
			return true
		}
	}
	return false
}

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

func (s *BuyerServiceDefault) Save(buyer *internal.Buyer) (err error) {
	all := s.repo.GetAll()
	ok := buyer.Parse()
	if !ok {
		err = BuyerUnprocessableEntity
		return
	}

	if cardNumberIdAlreadyInUse(buyer.CardNumberId, all) {
		err = CardNumberAlreadyInUse
		return
	}

	s.repo.Add(buyer)
	return
}

func (s *BuyerServiceDefault) Update(id int, buyerPatch internal.BuyerPatch) (err error) {
	all := s.repo.GetAll()
	_, ok := all[id]
	if !ok {
		err = BuyerNotFound
		return
	}

	if cardNumberIdAlreadyInUse(*buyerPatch.CardNumberId, all) {
		err = CardNumberAlreadyInUse
		return
	}

	s.repo.Update(id, buyerPatch)
	return
}

func (s *BuyerServiceDefault) Delete(id int) (err error) {
	all := s.repo.GetAll()
	_, ok := all[id]
	if !ok {
		err = BuyerNotFound
		return
	}

	s.repo.Delete(id)
	return
}

// ReportPurchaseOrders returns all purchase orders of all buyers
func (s *BuyerServiceDefault) ReportPurchaseOrders() (po []internal.PurchaseOrdersByBuyer, err error) {
	po, err = s.repo.ReportPurchaseOrders()
	// Check if there is no buyers records
	if len(po) == 0 {
		return nil, BuyerNotFound
	}
	return
}

// ReportPurchaseOrdersById returns all purchase orders of a specific buyer
func (s *BuyerServiceDefault) ReportPurchaseOrdersById(id int) (po []internal.PurchaseOrdersByBuyer, err error) {
	po, err = s.repo.ReportPurchaseOrdersById(id)
	// Check if there is no records for the given buyer
	if len(po) == 0 {
		return nil, BuyerNotFound
	}
	return
}
