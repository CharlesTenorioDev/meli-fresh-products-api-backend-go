package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

func NewProductService(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{repo: rp}
}

type ProductDefault struct {
	repo internal.ProductRepository
}

func (s *ProductDefault) GetAll() (v map[int]internal.Product, err error) {
	v, err = s.repo.GetAll()
	return
}
func (s *ProductDefault) GetByID(id int) (internal.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return internal.Product{}, err
	}
	return product, nil
}

func (s *ProductDefault) Create(product internal.Product) (internal.Product, error) {
	product, err := s.repo.Create(product)
	if err != nil {
		return internal.Product{}, err
	}
	return product, nil
}
func (s *ProductDefault) Update(product internal.Product) (internal.Product, error) {
	product, err := s.repo.Update(product)
	if err != nil {
		return internal.Product{}, err
	}
	return product, nil
}
func (s *ProductDefault) Delete(id int) (error) {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}