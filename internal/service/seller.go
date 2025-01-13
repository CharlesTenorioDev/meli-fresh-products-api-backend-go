package service

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type SellerServiceDefault struct {
	rp       internal.SellerRepository
	locality internal.LocalityRepository
}

func NewSellerServiceDefault(rp internal.SellerRepository, locality internal.LocalityRepository) *SellerServiceDefault {
	return &SellerServiceDefault{
		rp:       rp,
		locality: locality,
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

	if seller.CID == sellerCid.CID {
		return internal.ErrSellerCIDAlreadyExists
	}

	_, err = s.locality.FindByID(seller.Locality)
	if err != nil {
		return err
	}

	err = s.rp.Save(seller)
	if err != nil {
		return err
	}
	return nil
}

func (s *SellerServiceDefault) Update(id int, updatedSeller internal.SellerPatch) (internal.Seller, error) {
	actualSeller, err := s.FindByID(id)
	if err != nil {
		return internal.Seller{}, err
	}

	if updatedSeller.CID != nil {
		sellerCid, err := s.rp.FindByCID(*updatedSeller.CID)
		if err != nil && !errors.Is(err, internal.ErrSellerNotFound) {
			return internal.Seller{}, err
		}
		if *updatedSeller.CID == sellerCid.CID && actualSeller.ID != sellerCid.ID {
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

	if updatedSeller.Locality != nil {
		_, err := s.locality.FindByID(*updatedSeller.Locality)
		if err != nil {
			return internal.Seller{}, internal.ErrLocalityNotFound
		}
		actualSeller.Locality = *updatedSeller.Locality
	}

	err = s.rp.Update(&actualSeller)

	return actualSeller, err
}

func (s *SellerServiceDefault) Delete(id int) error {
	return s.rp.Delete(id)
}
