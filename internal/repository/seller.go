package repository

import (
	"encoding/json"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"log"
	"os"
)

type SellerRepoMap struct {
	db     map[int]internal.Seller
	lastId int
}

func NewSellerRepoMap() *SellerRepoMap {
	var sellers []internal.Seller
	db := make(map[int]internal.Seller)
	file, err := os.Open("db/sellers.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&sellers)
	if err != nil {
		log.Fatal(err)
	}

	var lastId int
	for i, s := range sellers {
		db[i+1] = s
		lastId++
	}

	return &SellerRepoMap{
		db:     db,
		lastId: lastId,
	}
}

func (s *SellerRepoMap) Save(seller *internal.Seller) error {
	id := s.lastId + 1

	_, ok := s.db[id]

	if ok {
		return internal.ErrSellerConflict
	}

	seller.ID = id
	s.db[id] = *seller
	s.lastId = id
	return nil
}

func (s *SellerRepoMap) FindByID(id int) (internal.Seller, error) {
	seller, ok := s.db[id]
	if !ok {
		return internal.Seller{}, internal.ErrSellerNotFound
	}
	return seller, nil
}

func (s *SellerRepoMap) FindByCID(cid int) (*internal.Seller, error) {
	sellers := s.db
	for _, seller := range sellers {
		if seller.CID == cid {
			return &seller, nil
		}
	}
	return nil, internal.ErrSellerNotFound
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
