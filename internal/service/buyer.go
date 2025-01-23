package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ErrBuyerNotFound                 = errors.New("buyer not found")
	ErrBuyerAlreadyExists            = errors.New("buyer already exists")
	ErrCardNumberAlreadyInUse        = errors.New("buyer with given card number already registered")
	ErrBuyerUnprocessableEntity      = errors.New("couldn't parse buyer")
	ErrPurchaseOrdersByBuyerNotFound = errors.New("purchase orders not found for the given buyer")
	ErrPurchaseOrdersNotFound        = errors.New("purchase orders not found for any buyer")
)

func cardNumberIDAlreadyInUse(cardNumber string, buyers map[int]internal.Buyer) bool {
	for _, b := range buyers {
		if b.CardNumberID == cardNumber {
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
		err = ErrBuyerNotFound
	}

	return
}

func (s *BuyerServiceDefault) Save(buyer *internal.Buyer) (err error) {
	all := s.repo.GetAll()

	ok := buyer.Parse()
	if !ok {
		err = ErrBuyerUnprocessableEntity
		return
	}

	if cardNumberIDAlreadyInUse(buyer.CardNumberID, all) {
		err = ErrCardNumberAlreadyInUse
		return
	}

	s.repo.Add(buyer)

	return
}

func (s *BuyerServiceDefault) Update(id int, buyerPatch internal.BuyerPatch) (err error) {
	all := s.repo.GetAll()

	_, ok := all[id]
	if !ok {
		err = ErrBuyerNotFound
		return
	}

	if cardNumberIDAlreadyInUse(*buyerPatch.CardNumberID, all) {
		err = ErrCardNumberAlreadyInUse
		return
	}

	s.repo.Update(id, buyerPatch)

	return
}

func (s *BuyerServiceDefault) Delete(id int) (err error) {
	all := s.repo.GetAll()

	_, ok := all[id]
	if !ok {
		err = ErrBuyerNotFound
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
		return nil, ErrPurchaseOrdersNotFound
	}
  
	return
}

// ReportPurchaseOrdersByID returns all purchase orders of a specific buyer
func (s *BuyerServiceDefault) ReportPurchaseOrdersByID(id int) (po []internal.PurchaseOrdersByBuyer, err error) {
	// Check if the buyer exists
	_, err = s.FindByID(id)
	if err != nil {
		return nil, err
	}
	// Get the purchase orders of the given buyer
	po, err = s.repo.ReportPurchaseOrdersByID(id)
	// Check if there is no records for the given buyer
	if len(po) == 0 {
		return nil, ErrPurchaseOrdersByBuyerNotFound
	}

	return
}
