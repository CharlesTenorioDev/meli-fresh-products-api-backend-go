package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewProductService(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{repo: rp}
}

type ProductDefault struct {
	repo internal.ProductRepository
}

func (s *ProductDefault) GetAll() (v map[int]internal.Product, err error) {
	v, err = s.repo.FindAll()
	return
}
func (s *ProductDefault) GetByID(id int) (internal.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return internal.Product{}, err
	}
	return product, nil
}

func (s *ProductDefault) Create(product internal.Product) (internal.Product, error) {
	existingProducts, err :=  s.repo.FindAll()
	if err != nil {
		return product, err
	}
	if err := ValidateProduct(product); err != nil {
		return product, err
	}
	if IsProductCodeExists(existingProducts, product.ProductCode) {
		return product, errors.New("product code already exists")
	}
	
	product.Id = GenerateNewID(existingProducts)
	product, err = s.repo.Save(product)
	if err != nil {
		return internal.Product{}, err
	}
	return product, nil
}
func (s *ProductDefault) Update(product internal.Product) (internal.Product, error) {
	existingProducts, err :=  s.repo.FindAll()
	if err != nil {
		return product, err
	}
	if err := ValidateProduct(product); err != nil {
		return product, err 
	}
	if IsProductCodeExists(existingProducts, product.ProductCode) {
		return product, errors.New("product code already exists")
	}
	productupdate, err := s.repo.Update(product)
	if err != nil {
		return internal.Product{}, err
	}
	return productupdate, nil
}
func (s *ProductDefault) Delete(id int) (error) {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func GenerateNewID(existingProducts map[int]internal.Product) int {
	maxID := 0
	for _, p := range existingProducts {
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	return maxID + 1
}
func IsProductCodeExists(existingProducts map[int]internal.Product, productCode string) bool {
	for _, p := range existingProducts {
		if p.ProductCode == productCode {
			return true
		}
	}
	return false
}
func ValidateProduct(product internal.Product) error {
	if product.ProductCode == "" || 
	product.Description == "" ||
	product.Height <= 0 || 
	product.Width <= 0 || 
	product.NetWeight <= 0 || 
	product.ExpirationRate.IsZero()|| 
	product.RecommendedFreezingTemperature < -273.15||
	product.FreezingRate < -273.15 ||
	product.ProductTypeId <= 0 ||
	product.SellerId <= 0 {
		return errors.New("all fields must be valid and filled")
	}
	return nil
}
