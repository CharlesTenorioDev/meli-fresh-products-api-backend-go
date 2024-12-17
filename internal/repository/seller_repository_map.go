package repository

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type SellerRepoMap struct {
	db map[int]internal.Seller
}

func NewSellerRepoMap(db map[int]internal.Seller) *SellerRepoMap {
	return &SellerRepoMap{db: db}
}

func (s *SellerRepoMap) Save(seller *internal.Seller) (int, error) {
	id := len(s.db) + 1

	_, ok := s.db[id]

	if ok {
		return 0, internal.ErrSellerConflict
	}

	seller.ID = id
	s.db[id] = *seller
	return id, nil
}

func (s *SellerRepoMap) FindByID(id int) (internal.Seller, error) {
	seller, ok := s.db[id]
	if !ok {
		return internal.Seller{}, internal.ErrSellerNotFound
	}
	return seller, nil
}

func (s *SellerRepoMap) FindByCID(cid int) (internal.Seller, error) {
	sellers := s.db
	for _, seller := range sellers {
		if seller.CID == cid {
			return seller, nil
		}
	}
	return internal.Seller{}, internal.ErrSellerNotFound
}

func (s *SellerRepoMap) FindAll() ([]internal.Seller, error) {
	var sellers []internal.Seller

	if len(s.db) == 0 {
		return nil, internal.ErrSellerNotFound
	}

	for _, seller := range s.db {
		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func (s *SellerRepoMap) Update(id int, seller *internal.Seller) error {
	_, ok := s.db[id]

	if !ok {
		return internal.ErrSellerNotFound
	}
	s.db[id] = *seller
	return nil
}

func (s *SellerRepoMap) Delete(id int) error {
	_, ok := s.db[id]
	if !ok {
		return internal.ErrSellerNotFound
	}
	delete(s.db, id)
	return nil
}
