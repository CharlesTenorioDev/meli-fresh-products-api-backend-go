package repository

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewProductMap() *ProductMap{
	var products []internal.Product
	db := make(map[int]*internal.Product)
	file, err := os.Open("db/product.json")

	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&products)
	if err != nil {
		log.Fatal(err)
	}
	for index, p := range products{
		db[index] = &p
	}
	return &ProductMap{db: db}
}

type ProductMap struct {
	db map[int]*internal.Product
}

func (r *ProductMap)GetAll() (db map[int]internal.Product, err error){
	db = make(map[int]internal.Product)
	for index, p := range r.db{
		db[index] = *p
	}
	return
}
func (r *ProductMap) GetByID(id int) (internal.Product, error) {
	product, exists := r.db[id]
	if !exists {
		return internal.Product{}, errors.New("product not found")
	}
	return *product, nil
}

func (r *ProductMap) Create(product internal.Product) (internal.Product, error) {
	if _, exists := r.db[product.Id]; exists {
		return internal.Product{}, errors.New("product with this ID already exists")
	}

	r.db[product.Id] = &product 

	var products []internal.Product
	file, err := os.Open("db/product.json")
	if err != nil {
		return product, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&products); err != nil {
		return product, err
	}

	products = append(products, product)

	file, err = os.Create("db/product.json") 
	if err != nil {
		return product, err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(products); err != nil {
		return product, err
	}

	return product, nil 
}



