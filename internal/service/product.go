package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

func NewProductService(rp internal.ProductRepository) *ProductDefault{
	return &ProductDefault{repo:rp}
}

type ProductDefault struct{
	repo internal.ProductRepository
}

func(s *ProductDefault) GetAll()(v map[int]internal.Product, err error){
	v, err = s.repo.GetAll()
	return
}