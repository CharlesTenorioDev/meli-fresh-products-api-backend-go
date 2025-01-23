package service_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryProductMock struct {
	mock.Mock
}

func (m *repositoryProductMock) FindAll() ([]internal.Product, error) {
	args := m.Called()
	return args.Get(0).([]internal.Product), args.Error(1)
}

func (r *repositoryProductMock) FindByID(id int) (internal.Product, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (r *repositoryProductMock) Save(product internal.Product) (internal.Product, error) {
	args := r.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (r *repositoryProductMock) Update(product internal.Product) (internal.Product, error) {
	args := r.Called(product)
	return args.Get(0).(internal.Product), args.Error(1)
}

func (r *repositoryProductMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func (m *repositoryProductMock) FindAllRecord() ([]internal.ProductRecordsJSONCount, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductRecordsJSONCount), args.Error(1)
}

func (r *repositoryProductMock) FindByIDRecord(id int) (internal.ProductRecordsJSONCount, error) {
	args := r.Called(id)
	return args.Get(0).(internal.ProductRecordsJSONCount), args.Error(1)
}

// Criação dos mocks dos repositórios

func TestProductServiceDefault_GetAll(t *testing.T) { //Se a lista tiver "n" elementos, ele retornará um número do total de elementos.
	t.Run("find_all", func(t *testing.T) {
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)

		productRepo.On("FindByID", 1).Return(expectedProduct, nil)

		// Chamada do método que será testado
		products, err := svc.GetByID(1)
		assert.Nil(t, err)
		assert.Equal(t, expectedProduct, products)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) { // Se o elemento pesquisado por id não existir, retornará null
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{product}, nil)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductCodeAlreadyExists, err)
	})

	t.Run("should error product seller not exists", func(t *testing.T) {
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerID).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", product.ProductTypeID).Return(internal.ProductType{}, internal.ErrProductTypeIdNotFound)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductTypeIdNotFound, err)
	})

	t.Run("should error if repository fails to save", func(t *testing.T) {
		//cria um erro de servidor
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

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

}

func TestProductServiceDefault_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {

		//Se a exclusão for bem-sucedida, o item não aparecerá na lista.
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("Delete", 1).Return(nil)
		err := svc.Delete(1)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("delete_non_existent", func(t *testing.T) {

		//Quando o produto não existir, será retornado null
		productRepo := new(repositoryProductMock)
		sellerRepo := new(repositoryMock)
		productTypeRepo := new(service.ProductTypeRepositoryMock)

		svc := service.NewProductService(productRepo, sellerRepo, productTypeRepo)
		productRepo.On("Delete", 1).Return(internal.ErrProductNotFound)
		err := svc.Delete(1)

		// Verifica se não houve erro
		assert.NotNil(t, err)
	})
}
