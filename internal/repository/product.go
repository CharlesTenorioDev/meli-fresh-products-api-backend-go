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

func (r *ProductMap) Update(product internal.Product) (internal.Product, error) {
	// Verifica se o produto existe
	if _, exists := r.db[product.Id]; !exists {
		return internal.Product{}, errors.New("product not found")
	}

	// Atualiza o produto no mapa
	r.db[product.Id] = &product 

	// Lê os produtos existentes do arquivo
	var products []internal.Product
	file, err := os.Open("db/product.json")
	if err != nil {
		return internal.Product{}, err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&products); err != nil {
		return internal.Product{}, err
	}

	// Atualiza o produto na lista existente
	for i, p := range products {
		if p.Id == product.Id {
			products[i] = product // Atualiza o produto na lista
			break
		}
	}

	// Salva a lista atualizada de produtos no arquivo JSON
	file, err = os.Create("db/product.json") 
	if err != nil {
		return internal.Product{}, err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(products); err != nil {
		return internal.Product{}, err
	}

	return product, nil // Retorna o produto atualizado
}

func (r *ProductMap) Delete(id int) error {
	// Verifica se o produto existe
	if _, exists := r.db[id]; !exists {
		return errors.New("product not found")
	}

	// Remove o produto do mapa
	delete(r.db, id)

	// Lê os produtos existentes do arquivo
	var products []internal.Product
	file, err := os.Open("db/product.json")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&products); err != nil {
		return err
	}

	// Cria uma nova lista sem o produto excluído
	var updatedProducts []internal.Product
	for _, p := range products {
		if p.Id != id {
			updatedProducts = append(updatedProducts, p)
		}
	}

	// Salva a lista atualizada de produtos no arquivo JSON
	file, err = os.Create("db/product.json")
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(updatedProducts); err != nil {
		return err
	}

	return nil // Retorna nil se a exclusão for bem-sucedida
}



