package repository

import (
	"encoding/json"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"log"
	"os"
)

type SellerMap struct {
	db     map[int]internal.Seller
	lastId int
}

func NewSellerMap() *SellerMap {
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

	for i, s := range sellers {
		db[i+1] = s
	}

	return &SellerMap{
		db:     db,
		lastId: len(sellers),
	}
}

func (s *SellerMap) Save(seller *internal.Seller) (err error) {
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

func (s *SellerMap) FindByID(id int) (seller internal.Seller, err error) {
	seller, ok := s.db[id]
	if !ok {
		return internal.Seller{}, internal.ErrSellerNotFound
	}
	return seller, nil
}

func (s *SellerMap) FindByCID(cid int) (seller internal.Seller, err error) {
	sellers := s.db
	for _, seller := range sellers {
		if seller.CID == cid {
			return seller, nil
		}
	}
	return internal.Seller{}, internal.ErrSellerNotFound
}

func (s *SellerMap) FindAll() (sellers []internal.Seller, err error) {
	if len(s.db) == 0 {
		return nil, internal.ErrSellerNotFound
	}

	for _, seller := range s.db {
		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func (s *SellerMap) Update(seller *internal.Seller) (err error) {
	_, ok := s.db[seller.ID]

	if !ok {
		return internal.ErrSellerNotFound
	}
	s.db[seller.ID] = *seller
	return nil
}

func (s *SellerMap) Delete(id int) (err error) {
	_, ok := s.db[id]
	if !ok {
		return internal.ErrSellerNotFound
	}
	delete(s.db, id)
	return nil
}