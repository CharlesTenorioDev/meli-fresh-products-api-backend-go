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

func (s *SellerRepoMap) Save(seller *internal.Seller) (id int, err error) {
	id = len(s.db) + 1

	_, ok := s.db[id]

	if ok {
		return 0, internal.ErrSellerConflict
	}

	seller.ID = id
	s.db[id] = *seller
	return
}

func (s *SellerRepoMap) Get(id int) (sl internal.Seller, err error) {
	sl, ok := s.db[id]
	if !ok {
		return internal.Seller{}, internal.ErrSellerNotFound
	}
	return
}

func (s *SellerRepoMap) GetAll() (sls []internal.Seller, err error) {
	var sellersList []internal.Seller

	if len(s.db) == 0 {
		return nil, internal.ErrSellerNotFound
	}

	for _, seller := range s.db {
		sellersList = append(sellersList, seller)
	}
	sls = sellersList

	return
}

func (s *SellerRepoMap) Update(id int, sl *internal.Seller) (err error) {
	_, ok := s.db[id]

	if !ok {
		return internal.ErrSellerNotFound
	}
	s.db[id] = *sl
	return
}

func (s *SellerRepoMap) Delete(id int) (err error) {
	_, ok := s.db[id]
	if !ok {
		return internal.ErrSellerNotFound
	}
	delete(s.db, id)
	return
}
