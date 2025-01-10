package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ProductUnprocessableEntity = errors.New("all fields must be valid and filled")
	ProductCodeAlreadyExists   = errors.New("product code already exists")
	ProductNotExists           = errors.New("Error ID doesn't exists")
	SellerNotExists            = errors.New("Error fetching seller")
	ProductTypeNotExists       = errors.New("Error fetching product type")
)

func NewProductService(prRepo internal.ProductRepository, slRepo internal.SellerRepository, ptRepo internal.ProductTypeRepository) *ProductDefault {
	return &ProductDefault{
		productRepo:     prRepo,
		sellerRepo:      slRepo,
		productTypeRepo: ptRepo,
	}
}

type ProductDefault struct {
	productRepo     internal.ProductRepository
	sellerRepo      internal.SellerRepository
	productTypeRepo internal.ProductTypeRepository
}

func (s *ProductDefault) GetAll() (v []internal.Product, err error) {
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
		return product, ProductCodeAlreadyExists
	}

	_, err = s.sellerRepo.FindByID(product.SellerId)
	if err != nil {
		return product, SellerNotExists
	}
	_, err = s.productTypeRepo.FindByID(product.ProductTypeId)
	if err != nil {
		return product, ProductTypeNotExists
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
	existingProducts, err := s.productRepo.FindAll()
	if err != nil {
		return product, err
	}

	existingProduct, err := s.productRepo.FindByID(product.Id)
	if err != nil {
		return product, ProductNotExists
	}

	if product.ProductCode == "" {
		product.ProductCode = existingProduct.ProductCode
	}
	if product.Description == "" {
		product.Description = existingProduct.Description
	}
	if product.Height == 0 {
		product.Height = existingProduct.Height
	}
	if product.Width == 0 {
		product.Width = existingProduct.Width
	}
	if product.NetWeight == 0 {
		product.NetWeight = existingProduct.NetWeight
	}
	if product.ExpirationRate.IsZero() {
		product.ExpirationRate = existingProduct.ExpirationRate
	}
	if product.RecommendedFreezingTemperature == 0 {
		product.RecommendedFreezingTemperature = existingProduct.RecommendedFreezingTemperature
	}
	if product.FreezingRate == 0 {
		product.FreezingRate = existingProduct.FreezingRate
	}
	if product.ProductTypeId == 0 {
		product.ProductTypeId = existingProduct.ProductTypeId
	}
	if product.SellerId == 0 {
		product.SellerId = existingProduct.SellerId
	}

	if IsProductCodeExists(existingProducts, product.ProductCode) {
		return product, ProductCodeAlreadyExists
	}

	_, err = s.sellerRepo.FindByID(product.SellerId)
	if err != nil {
		return product, SellerNotExists
	}
	_, err = s.productTypeRepo.FindByID(product.ProductTypeId)
	if err != nil {
		return product, ProductTypeNotExists
	}

	productUpdate, err := s.productRepo.Update(product)
	if err != nil {
		return internal.Product{}, err
	}

	return productUpdate, nil
}

func (s *ProductDefault) Delete(id int) error {
	err := s.productRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func GenerateNewID(existingProducts []internal.Product) int {
	maxID := 0
	for _, p := range existingProducts {
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	return maxID + 1
}
func IsProductCodeExists(existingProducts []internal.Product, productCode string) bool {
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
		product.ExpirationRate.IsZero() ||
		product.RecommendedFreezingTemperature < -273.15 ||
		product.FreezingRate < -273.15 ||
		product.ProductTypeId <= 0 ||
		product.SellerId <= 0 {
		return ProductUnprocessableEntity
	}
	return nil
}
