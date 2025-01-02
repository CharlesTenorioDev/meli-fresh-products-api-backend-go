package repository

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewProductMap() *ProductMap {
	var products []internal.Product
	db := make(map[int]*internal.Product)
	file, err := os.Open("db/product.json")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&products)
	if err != nil {
		log.Fatal(err)
	}

	for _, p := range products {
		db[p.Id] = &p
	}
	return &ProductMap{db: db}
}

type ProductMap struct {
	db map[int]*internal.Product
}

func (r *ProductMap) FindAll() (db []internal.Product, err error) {
	products := make([]internal.Product, 0, len(r.db))
	for _, p := range r.db {
		products = append(products, *p)
	}
	return products, nil
}
func (r *ProductMap) FindByID(id int) (internal.Product, error) {
	product, exists := r.db[id]
	if !exists {
		return internal.Product{}, errors.New("product not found")
	}
	return *product, nil
}

func (r *ProductMap) Save(product internal.Product) (internal.Product, error) {
	if _, exists := r.db[product.Id]; exists {
		return internal.Product{}, errors.New("product with this ID already exists")
	}

	r.db[product.Id] = &product
	return product, nil
}

func (r *ProductMap) Update(product internal.Product) (internal.Product, error) {
	if _, exists := r.db[product.Id]; !exists {
		return internal.Product{}, errors.New("product not found")
	}

	r.db[product.Id] = &product

	return product, nil
}

func (r *ProductMap) Delete(id int) error {
	if _, exists := r.db[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.db, id)

	return nil
}
