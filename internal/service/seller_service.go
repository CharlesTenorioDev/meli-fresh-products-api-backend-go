package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

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
	return s.rp.Save(seller)
}

func (s *SellerServiceDefault) Update(seller *internal.Seller) error {
	return s.rp.Update(seller)
}

func (s *SellerServiceDefault) Delete(id int) error {
	return s.rp.Delete(id)
}
