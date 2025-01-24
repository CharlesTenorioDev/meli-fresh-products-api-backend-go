package service_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func NewRepositoryProductMock() *RepositoryProductMock {
	return &RepositoryProductMock{}
}

type RepositoryProductMock struct {
	mock.Mock
}

func (m *RepositoryProductMock) FindAll() ([]internal.Product, error) {
	args := m.Called()
	return args.Get(0).([]internal.Product), args.Error(1)
}

func (r *RepositoryProductMock) FindByID(id int) (internal.Product, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (r *RepositoryProductMock) Save(product internal.Product) (internal.Product, error) {
	args := r.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (r *RepositoryProductMock) Update(product internal.Product) (internal.Product, error) {
	args := r.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (r *RepositoryProductMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func (m *RepositoryProductMock) FindAllRecord() ([]internal.ProductRecordsJSONCount, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductRecordsJSONCount), args.Error(1)
}

func (r *RepositoryProductMock) FindByIDRecord(id int) (internal.ProductRecordsJSONCount, error) {
	args := r.Called(id)
	return args.Get(0).(internal.ProductRecordsJSONCount), args.Error(1)
}

func TestProductServiceDefault_GetAll(t *testing.T) { //Se a lista tiver "n" elementos, ele retornará um número do total de elementos.
	t.Run("find_all", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		// Produtos esperados que o repositório deve retornar
		expectedProducts := []internal.Product{
			{ID: 1, ProductCode: "P001"},
			{ID: 2, ProductCode: "P002"},
		}

		// Configuração do mock para o método FindAll
		productRepo.On("FindAll").Return(expectedProducts, nil)

		// Chamada do método que será testado
		products, err := svc.GetAll()
		assert.Nil(t, err)
		assert.Equal(t, expectedProducts, products)
	})
}

func TestProductServiceDefault_GetByID(t *testing.T) {
	t.Run("find_by_id_existent", func(t *testing.T) { //Se o elemento pesquisado pelo id existir, ele retornará as informações do elemento solicitado
		expectedProduct := internal.Product{ID: 1}

		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)

		productRepo.On("FindByID", 1).Return(expectedProduct, nil)

		// Chamada do método que será testado
		products, err := svc.GetByID(1)
		assert.Nil(t, err)
		assert.Equal(t, expectedProduct, products)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) { // Se o elemento pesquisado por id não existir, retornará null
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindByID", 1).Return(internal.Product{}, internal.ErrProductNotFound)

		// Chamada do método que será testado
		products, err := svc.GetByID(1)
		assert.NotNil(t, err)
		assert.Equal(t, internal.Product{}, products)
	})
}

func TestProductServiceDefault_Create(t *testing.T) {
	product := internal.Product{
		ID:                             1,
		ProductCode:                    "code 1",
		Description:                    "Example Product",
		Height:                         1,
		Length:                         1,
		NetWeight:                      1,
		ExpirationRate:                 1,
		RecommendedFreezingTemperature: 1,
		Width:                          1,
		FreezingRate:                   1,
		ProductTypeID:                  1,
		SellerID:                       1,
	}

	t.Run("create_ok", func(t *testing.T) { //Se ele contiver os campos necessários, será criado
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		// Configura o mock para as chamadas necessárias
		productRepo.On("FindAll").Return([]internal.Product{}, nil)                               // Configuração para FindAll
		productRepo.On("Save", product).Return(product, nil)                                      // Configuração para Save
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)                // Configuração para FindByID no sellerRepo
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, nil) // Configuração para FindByID no productTypeRepo

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("create_conflito", func(t *testing.T) { // Se o product_code já existir, ele não poderá ser criado.
		// Cria um product com code já existente.
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{product}, nil)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductCodeAlreadyExists, err)
	})

	t.Run("should error product seller not exists", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		// Cria um product com seller que não existe
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, internal.ErrSellerIdNotFound)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrSellerIdNotFound, err)
	})

	t.Run("should error product product-type not exists", func(t *testing.T) {
		// Cria um product com product-type não existe
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, internal.ErrProductTypeIDNotFound)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductTypeIDNotFound, err)
	})

	t.Run("should error if repository fails to save", func(t *testing.T) {
		//cria um erro de servidor
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, errors.New("repository error"))
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, nil)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, "repository error", err.Error())
	})
	t.Run("should error if repository FindAll", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, errors.New("repository error"))
		productRepo.On("Save", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, nil)
		_, err := svc.Create(product)

		assert.NotNil(t, err)
		assert.Equal(t, "repository error", err.Error())
	})
	t.Run("should error if product validation fails", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)

		productRepo.On("FindAll").Return([]internal.Product{}, nil)

		product := internal.Product{
			ID:          1,
			ProductCode: "",
			Description: "Example Product",
			Height:      1,
		}

		_, err := svc.Create(product)

		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductUnprocessableEntity, err)
	})

}

func TestProductServiceDefault_Update(t *testing.T) {
	product := internal.Product{
		ID:                             1,
		ProductCode:                    "code 1",
		Description:                    "Example Product",
		Height:                         1,
		Length:                         1,
		NetWeight:                      1,
		ExpirationRate:                 1,
		RecommendedFreezingTemperature: 1,
		Width:                          1,
		FreezingRate:                   1,
		ProductTypeID:                  1,
		SellerID:                       1,
	}

	t.Run("update_existent", func(t *testing.T) {
		//Quando a atualização de dados for bem-sucedida, o produto será devolvido com as informações atualizadas.

		// Configura o mock para as chamadas necessárias
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("FindByID", product.ID).Return(product, nil)                               // Configuração para FindAll
		productRepo.On("Update", product).Return(product, nil)                                    // Configuração para Save
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)                // Configuração para FindByID no sellerRepo
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, nil) // Configuração para FindByID no productTypeRepo

		// Executa o método que será testado
		_, err := svc.Update(product)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("update_non_existent", func(t *testing.T) {
		//Se o produto a ser atualizado não existir, será retornado null.

		// Configura o mock para as chamadas necessárias
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("FindByID", product.ID).Return(product, internal.ErrProductNotFound)       // Configuração para FindAll
		productRepo.On("Update", product).Return(product, nil)                                    // Configuração para Save
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)                // Configuração para FindByID no sellerRepo
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, nil) // Configuração para FindByID no productTypeRepo

		// Executa o método que será testado
		_, err := svc.Update(product)

		// Verifica se não houve erro
		assert.NotNil(t, err)
	})
	t.Run("should error if repository FindAll", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, errors.New("repository error"))
		productRepo.On("Update", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, nil)
		_, err := svc.Update(product)

		assert.NotNil(t, err)
		assert.Equal(t, "repository error", err.Error())
	})
	t.Run("should fill missing fields from existing product", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)

		existingProduct := internal.Product{
			ID:                             1,
			ProductCode:                    "code 1",
			Description:                    "Existing Description",
			Height:                         100,
			Width:                          50,
			Length:                         30,
			NetWeight:                      20,
			ExpirationRate:                 5,
			RecommendedFreezingTemperature: -18,
			FreezingRate:                   3,
			ProductTypeID:                  2,
			SellerID:                       3,
		}

		// Configuração do mock
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("FindByID", existingProduct.ID).Return(existingProduct, nil)
		productRepo.On("Update", mock.Anything).Return(existingProduct, nil)
		sellerRepo.On("FindByID", existingProduct.SellerID).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", existingProduct.ProductTypeID).Return(internal.ProductType{}, nil)

		// Produto de entrada com campos faltando ou zerados
		product := internal.Product{
			ID: 1, // Mesmo ID para garantir que os valores serão preenchidos
			// Os demais campos estão vazios ou zerados
		}

		updatedProduct, err := svc.Update(product)

		assert.Nil(t, err)

		// Verifica se os campos foram preenchidos corretamente a partir do existingProduct
		assert.Equal(t, existingProduct.ProductCode, updatedProduct.ProductCode)
		assert.Equal(t, existingProduct.Description, updatedProduct.Description)
		assert.Equal(t, existingProduct.Height, updatedProduct.Height)
		assert.Equal(t, existingProduct.Width, updatedProduct.Width)
		assert.Equal(t, existingProduct.Length, updatedProduct.Length)
		assert.Equal(t, existingProduct.NetWeight, updatedProduct.NetWeight)
		assert.Equal(t, existingProduct.ExpirationRate, updatedProduct.ExpirationRate)
		assert.Equal(t, existingProduct.RecommendedFreezingTemperature, updatedProduct.RecommendedFreezingTemperature)
		assert.Equal(t, existingProduct.FreezingRate, updatedProduct.FreezingRate)
		assert.Equal(t, existingProduct.ProductTypeID, updatedProduct.ProductTypeID)
		assert.Equal(t, existingProduct.SellerID, updatedProduct.SellerID)
	})

}

func TestProductServiceDefault_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {

		//Se a exclusão for bem-sucedida, o item não aparecerá na lista.
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("Delete", 1).Return(nil)
		err := svc.Delete(1)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("delete_non_existent", func(t *testing.T) {

		//Quando o produto não existir, será retornado null
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("Delete", 1).Return(internal.ErrProductNotFound)
		err := svc.Delete(1)

		// Verifica se não houve erro
		assert.NotNil(t, err)
	})
}

func TestProductServiceDefault_GetByIDRecord(t *testing.T) {
	t.Run("find_by_id_existent", func(t *testing.T) { //Se o elemento pesquisado pelo id existir, ele retornará as informações do elemento solicitado
		expectedProduct := internal.ProductRecordsJSONCount{ProductID: 1}

		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)

		productRepo.On("FindByIDRecord", 1).Return(expectedProduct, nil)

		// Chamada do método que será testado
		products, err := svc.GetByIDRecord(1)
		assert.Nil(t, err)
		assert.Equal(t, expectedProduct, products)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindByIDRecord", 1).Return(internal.ProductRecordsJSONCount{}, internal.ErrProductNotFound)

		// Chamada do método que será testado
		products, err := svc.GetByIDRecord(1)
		assert.NotNil(t, err)
		assert.Equal(t, internal.ProductRecordsJSONCount{}, products)
	})
}

func TestProductServiceDefault_GetAllRecord(t *testing.T) {
	t.Run("FindAllRecord", func(t *testing.T) {
		productRepo := new(RepositoryProductMock)
		sellerRepo := new(sellerRepositoryMock)
		productTypeRepo := new(ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		// Produtos esperados que o repositório deve retornar
		expectedProducts := []internal.ProductRecordsJSONCount{
			{ProductID: 1, Description: "P001", RecordsCount: 1},
			{ProductID: 2, Description: "P002", RecordsCount: 1},
		}

		// Configuração do mock para o método FindAll
		productRepo.On("FindAllRecord").Return(expectedProducts, nil)

		// Chamada do método que será testado
		products, err := svc.GetAllRecord()
		assert.Nil(t, err)
		assert.Equal(t, expectedProducts, products)
	})
}

func TestValidateProduct(t *testing.T) {
	validProduct := internal.Product{
		ProductCode:                    "P12345",
		Description:                    "Sample Product",
		Height:                         10.5,
		Length:                         20.0,
		Width:                          15.0,
		NetWeight:                      5.0,
		ExpirationRate:                 365.0,
		RecommendedFreezingTemperature: -20.0,
		FreezingRate:                   -10.0,
		ProductTypeID:                  1,
		SellerID:                       100,
	}

	tests := []struct {
		name    string
		product internal.Product
		wantErr error
	}{
		{
			name:    "Valid Product",
			product: validProduct,
			wantErr: nil,
		},
		{
			name: "Empty ProductCode",
			product: func() internal.Product {
				p := validProduct
				p.ProductCode = ""
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Empty Description",
			product: func() internal.Product {
				p := validProduct
				p.Description = ""
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Negative Height",
			product: func() internal.Product {
				p := validProduct
				p.Height = -1.0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Zero Length",
			product: func() internal.Product {
				p := validProduct
				p.Length = 0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Negative Width",
			product: func() internal.Product {
				p := validProduct
				p.Width = -5.0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Negative NetWeight",
			product: func() internal.Product {
				p := validProduct
				p.NetWeight = -2.5
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Zero ExpirationRate",
			product: func() internal.Product {
				p := validProduct
				p.ExpirationRate = 0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Below Absolute Zero RecommendedFreezingTemperature",
			product: func() internal.Product {
				p := validProduct
				p.RecommendedFreezingTemperature = -300.0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Below Absolute Zero FreezingRate",
			product: func() internal.Product {
				p := validProduct
				p.FreezingRate = -300.0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Zero ProductTypeID",
			product: func() internal.Product {
				p := validProduct
				p.ProductTypeID = 0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Negative SellerID",
			product: func() internal.Product {
				p := validProduct
				p.SellerID = -10
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
		{
			name: "Multiple Invalid Fields",
			product: func() internal.Product {
				p := validProduct
				p.ProductCode = ""
				p.Height = -5.0
				p.SellerID = 0
				return p
			}(),
			wantErr: internal.ErrProductUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateProduct(tt.product)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestIsProductCodeExists(t *testing.T) {
	tests := []struct {
		name             string
		existingProducts []internal.Product
		productCode      string
		expectedExists   bool
	}{
		{
			name: "Product code exists",
			existingProducts: []internal.Product{
				{ProductCode: "P001"},
				{ProductCode: "P002"},
			},
			productCode:    "P001",
			expectedExists: true,
		},
		{
			name: "Product code does not exist",
			existingProducts: []internal.Product{
				{ProductCode: "P001"},
				{ProductCode: "P002"},
			},
			productCode:    "P003",
			expectedExists: false,
		},
		{
			name:             "Empty product list",
			existingProducts: []internal.Product{},
			productCode:      "P001",
			expectedExists:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsProductCodeExists(tt.existingProducts, tt.productCode)
			assert.Equal(t, tt.expectedExists, result)
		})
	}
}

func TestGenerateNewID(t *testing.T) {
	tests := []struct {
		name             string
		existingProducts []internal.Product
		expectedID       int
	}{
		{
			name:             "Empty product list",
			existingProducts: []internal.Product{},
			expectedID:       1,
		},
		{
			name: "Single product",
			existingProducts: []internal.Product{
				{ID: 1},
			},
			expectedID: 2,
		},
		{
			name: "Multiple products with sequential IDs",
			existingProducts: []internal.Product{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			},
			expectedID: 4,
		},
		{
			name: "Multiple products with non-sequential IDs",
			existingProducts: []internal.Product{
				{ID: 5},
				{ID: 10},
				{ID: 3},
			},
			expectedID: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GenerateNewID(tt.existingProducts)

			if result != tt.expectedID {
				assert.Equal(t, tt.expectedID, result)
			}
		})
	}
}
