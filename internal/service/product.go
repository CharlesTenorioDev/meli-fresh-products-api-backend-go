package service

import (
	"errors"
	"fmt"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewProductService(prRepo internal.ProductRepository, slRepo internal.SellerRepository, ptRepo internal.ProductTypeRepository) *ProductDefault {
	return &ProductDefault{
		productRepo: prRepo, 
		sellerRepo: slRepo, 
		productTypeRepo: ptRepo,

	}
}
type ProductDefault struct {
	productRepo internal.ProductRepository
	sellerRepo internal.SellerRepository
	productTypeRepo internal.ProductTypeRepository
}

func (s *ProductDefault) GetAll() (v map[int]internal.Product, err error) {
	v, err = s.productRepo.FindAll()
	return
}
func (s *ProductDefault) GetByID(id int) (internal.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return internal.Product{}, err
	}
	return product, nil
}

func (s *ProductDefault) Create(product internal.Product) (internal.Product, error) {
    existingProducts, err := s.productRepo.FindAll()
    if err != nil {
        return product, err
    }

    if err := ValidateProduct(product); err != nil {
        return product, err
    }

    if IsProductCodeExists(existingProducts, product.ProductCode) {
        return product, errors.New("product code already exists")
    }

    _, err = s.sellerRepo.FindByID(product.SellerId)
    if err != nil {
        fmt.Println("Error fetching seller:", err)
        return product, err
    }
	_, err = s.productTypeRepo.FindByID(product.ProductTypeId)
    if err != nil {
        fmt.Println("Error fetching Product type:", err)
        return product, err
    }
    // Gera um novo ID para o produto
    product.Id = GenerateNewID(existingProducts)

    // Salva o novo produto no repositório
    product, err = s.productRepo.Save(product)
    if err != nil {
        return internal.Product{}, err
    }
    
    return product, nil
}

func (s *ProductDefault) Update(product internal.Product) (internal.Product, error) {
	existingProducts, err :=  s.productRepo.FindAll()
	if err != nil {
		return product, err
	}
	if err := ValidateProduct(product); err != nil {
		return product, err 
	}
	_, err = s.sellerRepo.FindByID(product.SellerId)
    if err != nil {
        fmt.Println("Error fetching seller:", err)
        return product, err
    }
	_, err = s.productTypeRepo.FindByID(product.ProductTypeId)
    if err != nil {
        fmt.Println("Error fetching Product type:", err)
        return product, err
    }
	
	if IsProductCodeExists(existingProducts, product.ProductCode) {
		return product, errors.New("product code already exists")
	}
	productupdate, err := s.productRepo.Update(product)
	if err != nil {
		return internal.Product{}, err
	}
	return productupdate, nil
}
func (s *ProductDefault) Delete(id int) (error) {
	err := s.productRepo.Delete(id)
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
