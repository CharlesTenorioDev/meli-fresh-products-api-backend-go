package service

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type SellerServiceDefault struct {
	rp internal.SellerRepository
}

func NewSellerServiceDefault(rp internal.SellerRepository) *SellerServiceDefault {
	return &SellerServiceDefault{
		rp: rp,
	}
}

func (s *SellerServiceDefault) FindAll() ([]internal.Seller, error) {
	sellers, err := s.rp.FindAll()
	if err != nil {
		return nil, err
	}

	return sellers, nil
}

func (s *SellerServiceDefault) FindByID(id int) (internal.Seller, error) {
	seller, err := s.rp.FindByID(id)
	if err != nil {
		return internal.Seller{}, err
	}

	return seller, nil
}

func (s *SellerServiceDefault) Save(seller *internal.Seller) (int, error) {

	sellerCid, err := s.rp.FindByCID(seller.CID)
	if err != nil && !errors.Is(err, internal.ErrSellerNotFound) {
		return 0, err
	}

	if sellerCid != nil {
		return 0, internal.ErrSellerCIDAlreadyExists
	}

	id, err := s.rp.Save(seller)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *SellerServiceDefault) Update(seller *internal.Seller) error {
	return s.rp.Update(seller.ID, seller)
}

func (s *SellerServiceDefault) Delete(id int) error {
	return s.rp.Delete(id)
}
