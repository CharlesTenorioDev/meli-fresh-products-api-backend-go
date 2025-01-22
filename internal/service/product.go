package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
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
		return product, internal.ErrProductCodeAlreadyExists
	}

	_, err = s.sellerRepo.FindByID(product.SellerId)
	if err != nil {
		return product, internal.ErrSellerIdNotFound
	}
	_, err = s.productTypeRepo.FindByID(product.ProductTypeId)
	if err != nil {
		return product, internal.ErrProductTypeIdNotFound
	}

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
		return product, internal.ErrProductNotFound
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
	if product.Length == 0 {
		product.Length = existingProduct.Length
	}
	if product.NetWeight == 0 {
		product.NetWeight = existingProduct.NetWeight
	}
	if product.ExpirationRate == 0 {
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
		return product, internal.ErrProductCodeAlreadyExists
	}

	_, err = s.sellerRepo.FindByID(product.SellerId)
	if err != nil {
		return product, internal.ErrSellerIdNotFound
	}
	_, err = s.productTypeRepo.FindByID(product.ProductTypeId)
	if err != nil {
		return product, internal.ErrProductTypeNotFound
	}

	_, err = s.productRepo.Update(product)
	if err != nil {
		return internal.Product{}, err
	}

	return product, nil
}

func (s *ProductDefault) Delete(id int) error {
	err := s.productRepo.Delete(id)
	if err != nil {
		print(err)
		return err
	}
	return nil
}

func (s *ProductDefault) GetAllRecord() (v []internal.ProductRecordsJsonCount, err error) {
	v, err = s.productRepo.FindAllRecord()
	return
}

func (s *ProductDefault) GetByIdRecord(id int) (internal.ProductRecordsJsonCount, error) {
	product, err := s.productRepo.FindByIdRecord(id)
	if err != nil {
		return internal.ProductRecordsJsonCount{}, err
	}
	return product, nil
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
	if product.ProductCode == "" {
		return errors.New("ProductCode is empty.")
	}
	if product.Description == "" {
		return errors.New("Description is empty.")
	}
	if product.Height <= 0 {
		return errors.New("Height is empty.")
	}
	if product.Length <= 0 {
		return errors.New("Length is empty.")
	}
	if product.Width <= 0 {
		return errors.New("Width is empty.")
	}
	if product.NetWeight <= 0 {
		return errors.New("NetWeight is empty.")
	}
	if product.ExpirationRate <= 0 {
		return errors.New("ExpirationRate is empty.")
	}
	if product.RecommendedFreezingTemperature < -273.15 {
		return errors.New("RecommendedFreezingTemperature is empty.")
	}
	if product.FreezingRate < -273.15 {
		return errors.New("FreezingRate is empty.")
	}
	if product.ProductTypeId <= 0 {
		return errors.New("ProductTypeId is empty.")
	}
	if product.SellerId <= 0 {
		return errors.New("SellerId is empty.")
	}
	return nil
}
