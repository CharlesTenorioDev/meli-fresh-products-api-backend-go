package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type BuyerMap struct {
	db map[int]*internal.Buyer
}

func NewBuyerMap(dbPath string) *BuyerMap {
	var buyers []internal.Buyer
	db := make(map[int]*internal.Buyer)
	file, err := os.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&buyers)
	if err != nil {
		log.Fatal(err)
	}

	for i, b := range buyers {
		db[i] = &b
	}
	return &BuyerMap{
		db: db,
	}
}

func (r *BuyerMap) GetAll() (db map[int]internal.Buyer) {
	db = make(map[int]internal.Buyer)

	for i, b := range r.db {
		db[i] = *b
	}
	return
}

func (r *BuyerMap) Add(buyer *internal.Buyer) {
	id := len(r.db)
	buyer.ID = id
	r.db[id] = buyer
}

func (r *BuyerMap) Update(id int, buyer internal.BuyerPatch) {
	buyerToPatch := r.db[id]
	buyer.Patch(buyerToPatch)
}

func (r *BuyerMap) Delete(id int) {
	delete(r.db, id)
}
