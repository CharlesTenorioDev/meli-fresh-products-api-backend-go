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

func NewBuyerMap() *BuyerMap {
	var buyers []internal.Buyer
	db := make(map[int]*internal.Buyer)
	file, err := os.Open("db/buyer.json")
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

func (r *BuyerMap) AddProduct(id int, buyer internal.Buyer) {
	r.db[id] = &buyer
}

func (r *BuyerMap) UpdateBuyer(id int, buyer internal.BuyerPatch) {
	buyerToPatch := r.db[id]
	buyer.Patch(buyerToPatch)
}
