package repository

import (
	"encoding/json"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type BuyerMap struct {
	db     map[int]*internal.Buyer
	lastID int
}

func NewBuyerMap(dbPath string) (mp *BuyerMap, err error) {
	var buyers []internal.Buyer

	db := make(map[int]*internal.Buyer)

	file, err := os.Open(dbPath)
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&buyers)
	if err != nil {
		return
	}

	for i, b := range buyers {
		db[i] = &b
	}

	mp = &BuyerMap{
		db:     db,
		lastID: len(buyers),
	}
	return
}

func (r *BuyerMap) GetAll() (db map[int]internal.Buyer) {
	db = make(map[int]internal.Buyer)

	for i, b := range r.db {
		db[i] = *b
	}

	return
}

func (r *BuyerMap) Add(buyer *internal.Buyer) {
	id := r.lastID
	buyer.ID = id
	r.db[id] = buyer
	r.lastID++
}

func (r *BuyerMap) Update(id int, buyer internal.BuyerPatch) {
	buyerToPatch := r.db[id]
	buyer.Patch(buyerToPatch)
}

func (r *BuyerMap) Delete(id int) {
	delete(r.db, id)
}

func (r *BuyerMap) ReportPurchaseOrders() (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	return
}

func (r *BuyerMap) ReportPurchaseOrdersByID(id int) (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	return
}
