package service

import (
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

	_, err = s.sellerRepo.FindByID(product.SellerID)
	if err != nil {
		return product, internal.ErrSellerIdNotFound
	}

	_, err = s.productTypeRepo.FindByID(product.ProductTypeID)
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

	existingProduct, err := s.productRepo.FindByID(product.ID)

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

	if product.ProductTypeID == 0 {
		product.ProductTypeID = existingProduct.ProductTypeID
	}

	if product.SellerID == 0 {
		product.SellerID = existingProduct.SellerID
	}

	if IsProductCodeExists(existingProducts, product.ProductCode) {
		return product, internal.ErrProductCodeAlreadyExists
	}

	_, err = s.sellerRepo.FindByID(product.SellerID)
	if err != nil {
		return product, internal.ErrSellerIdNotFound
	}

	_, err = s.productTypeRepo.FindByID(product.ProductTypeID)
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

func (s *ProductDefault) GetAllRecord() (v []internal.ProductRecordsJSONCount, err error) {
	v, err = s.productRepo.FindAllRecord()
	return
}

func (s *ProductDefault) GetByIDRecord(id int) (internal.ProductRecordsJSONCount, error) {
	product, err := s.productRepo.FindByIDRecord(id)
	if err != nil {
		return internal.ProductRecordsJSONCount{}, err
	}

	return product, nil
}
func GenerateNewID(existingProducts []internal.Product) int {
	maxID := 0
	for _, p := range existingProducts {
		if p.ID > maxID {
			maxID = p.ID
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
	if product.ProductCode == "" || product.Description == "" || product.Height <= 0 || product.Length <= 0 ||
		product.Width <= 0 || product.NetWeight <= 0 || product.ExpirationRate <= 0 || product.RecommendedFreezingTemperature < -273.15 ||
		product.FreezingRate < -273.15 || product.ProductTypeID <= 0 || product.SellerID <= 0 {
		return internal.ErrProductUnprocessableEntity
	}

	return nil
}
