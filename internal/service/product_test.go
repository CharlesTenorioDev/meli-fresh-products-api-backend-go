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

func (m *repositoryProductMock) FindAllRecord() ([]internal.ProductRecordsJsonCount, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductRecordsJsonCount), args.Error(1)
}

func (r *repositoryProductMock) FindByIdRecord(id int) (internal.ProductRecordsJsonCount, error) {
	args := r.Called(id)
	return args.Get(0).(internal.ProductRecordsJsonCount), args.Error(1)
}

// Criação dos mocks dos repositórios

var productRepo = new(repositoryProductMock)
var sellerRepo = new(repositoryMock)
var productTypeRepo = new(service.ProductTypeRepositoryMock)

var svc = service.NewProductService(productRepo, sellerRepo, productTypeRepo)

func TestProductServiceDefault_GetAll(t *testing.T) { //Se a lista tiver "n" elementos, ele retornará um número do total de elementos.
	t.Run("find_all", func(t *testing.T) {
		// Produtos esperados que o repositório deve retornar
		expectedProducts := []internal.Product{
			{Id: 1, ProductCode: "P001"},
			{Id: 2, ProductCode: "P002"},
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
		expectedProduct := internal.Product{Id: 1}

		productRepo.On("FindByID", 1).Return(expectedProduct, nil)

		// Chamada do método que será testado
		products, err := svc.GetByID(1)
		assert.Nil(t, err)
		assert.Equal(t, expectedProduct, products)
	})

	t.Run("find_by_id_non_existent", func(t *testing.T) { // Se o elemento pesquisado por id não existir, retornará null

		productRepo.On("FindByID", 1).Return(internal.Product{}, internal.ErrProductNotFound)

		// Chamada do método que será testado
		products, err := svc.GetByID(1)
		assert.NotNil(t, err)
		assert.Equal(t, internal.Product{}, products)
	})
}

func TestProductServiceDefault_Create(t *testing.T) {
	product := internal.Product{
		Id:                             1,
		ProductCode:                    "code 1",
		Description:                    "Example Product",
		Height:                         1,
		Length:                         1,
		NetWeight:                      1,
		ExpirationRate:                 1,
		RecommendedFreezingTemperature: 1,
		Width:                          1,
		FreezingRate:                   1,
		ProductTypeId:                  1,
		SellerId:                       1,
	}

	t.Run("create_ok", func(t *testing.T) { //Se ele contiver os campos necessários, será criado

		// Configura o mock para as chamadas necessárias
		productRepo.On("FindAll").Return([]internal.Product{}, nil)                               // Configuração para FindAll
		productRepo.On("Save", product).Return(product, nil)                                      // Configuração para Save
		sellerRepo.On("FindByID", product.SellerId).Return(internal.Seller{}, nil)                // Configuração para FindByID no sellerRepo
		productTypeRepo.On("FindByID", product.ProductTypeId).Return(internal.ProductType{}, nil) // Configuração para FindByID no productTypeRepo

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("create_conflito", func(t *testing.T) { // Se o product_code já existir, ele não poderá ser criado.
		// Cria um product com code já existente.
		productRepo.On("FindAll").Return([]internal.Product{product}, nil)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductCodeAlreadyExists, err)
	})

	t.Run("should error product seller not exists", func(t *testing.T) {
		// Cria um product com seller que não existe
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerId).Return(internal.Seller{}, internal.ErrSellerIdNotFound)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrSellerIdNotFound, err)
	})

	t.Run("should error product product-type not exists", func(t *testing.T) {
		// Cria um product com product-type não existe
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, nil)
		sellerRepo.On("FindByID", product.SellerId).Return(internal.Seller{}, nil)
		sellerRepo.On("FindByID", product.ProductTypeId).Return(internal.ProductType{}, internal.ErrProductTypeIdNotFound)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductIdNotFound, err)
	})

	t.Run("should error if repository fails to save", func(t *testing.T) {
		//cria um erro de servidor
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("Save", product).Return(product, errors.New("repository error"))
		sellerRepo.On("FindByID", product.SellerId).Return(internal.Seller{}, nil)
		productTypeRepo.On("FindByID", product.ProductTypeId).Return(internal.ProductType{}, nil)

		// Executa o método que será testado
		_, err := svc.Create(product)

		// Verifica se o erro retornado é o esperado
		assert.NotNil(t, err)
		assert.Equal(t, "repository error", err.Error())
	})

}

func TestProductServiceDefault_Update(t *testing.T) {
	product := internal.Product{
		Id:                             1,
		ProductCode:                    "code 1",
		Description:                    "Example Product",
		Height:                         1,
		Length:                         1,
		NetWeight:                      1,
		ExpirationRate:                 1,
		RecommendedFreezingTemperature: 1,
		Width:                          1,
		FreezingRate:                   1,
		ProductTypeId:                  1,
		SellerId:                       1,
	}

	t.Run("update_existent", func(t *testing.T) {
		//Quando a atualização de dados for bem-sucedida, o produto será devolvido com as informações atualizadas.

		// Configura o mock para as chamadas necessárias
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("FindByID", product.Id).Return(product, nil)                               // Configuração para FindAll
		productRepo.On("Update", product).Return(product, nil)                                    // Configuração para Save
		sellerRepo.On("FindByID", product.SellerId).Return(internal.Seller{}, nil)                // Configuração para FindByID no sellerRepo
		productTypeRepo.On("FindByID", product.ProductTypeId).Return(internal.ProductType{}, nil) // Configuração para FindByID no productTypeRepo

		// Executa o método que será testado
		_, err := svc.Update(product)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("update_non_existent", func(t *testing.T) {
		//Se o produto a ser atualizado não existir, será retornado null.

		// Configura o mock para as chamadas necessárias
		productRepo.On("FindAll").Return([]internal.Product{}, nil)
		productRepo.On("FindByID", product.Id).Return(product, internal.ErrProductNotFound)       // Configuração para FindAll
		productRepo.On("Update", product).Return(product, nil)                                    // Configuração para Save
		sellerRepo.On("FindByID", product.SellerId).Return(internal.Seller{}, nil)                // Configuração para FindByID no sellerRepo
		productTypeRepo.On("FindByID", product.ProductTypeId).Return(internal.ProductType{}, nil) // Configuração para FindByID no productTypeRepo

		// Executa o método que será testado
		_, err := svc.Update(product)

		// Verifica se não houve erro
		assert.NotNil(t, err)
	})

}

func TestProductServiceDefault_Delete(t *testing.T) {
	t.Run("delete_ok", func(t *testing.T) {

		//Se a exclusão for bem-sucedida, o item não aparecerá na lista.
		productRepo.On("Delete", 1).Return(nil)
		err := svc.Delete(1)

		// Verifica se não houve erro
		assert.Nil(t, err)
	})

	t.Run("delete_non_existent", func(t *testing.T) {

		//Quando o produto não existir, será retornado null
		productRepo.On("Delete", 1).Return(internal.ErrProductNotFound)
		err := svc.Delete(1)

		// Verifica se não houve erro
		assert.NotNil(t, err)
	})
}
