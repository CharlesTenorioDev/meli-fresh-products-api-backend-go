package service

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/dto"
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

func (s *SellerServiceDefault) Save(seller *internal.Seller) error {

	sellerCid, err := s.rp.FindByCID(seller.CID)
	if err != nil && !errors.Is(err, internal.ErrSellerNotFound) {
		return err
	}

	if sellerCid != nil {
		return internal.ErrSellerCIDAlreadyExists
	}

	err = s.rp.Save(seller)
	if err != nil {
		return err
	}
	return nil
}

func (s *SellerServiceDefault) Update(id int, updatedSeller dto.SellersUpdateDto) (internal.Seller, error) {
	actualSeller, err := s.FindByID(id)
	if err != nil {
		return internal.Seller{}, err
	}

	if updatedSeller.CID != nil {
		sellerCid, err := s.rp.FindByCID(*updatedSeller.CID)
		if err != nil && !errors.Is(err, internal.ErrSellerNotFound) {
			return internal.Seller{}, err
		}
		if sellerCid != nil && actualSeller.ID != sellerCid.ID {
			return internal.Seller{}, internal.ErrSellerCIDAlreadyExists
		}
		actualSeller.CID = *updatedSeller.CID
	}

	if updatedSeller.CompanyName != nil {
		actualSeller.CompanyName = *updatedSeller.CompanyName
	}

	if updatedSeller.Address != nil {
		actualSeller.Address = *updatedSeller.Address
	}

	if updatedSeller.Telephone != nil {
		actualSeller.Telephone = *updatedSeller.Telephone
	}

	err = s.rp.Update(actualSeller.ID, &actualSeller)

	return actualSeller, err
}

func (s *SellerServiceDefault) Delete(id int) error {
	return s.rp.Delete(id)
}
