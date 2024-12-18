package repository

import "github.com/meli-fresh-products-api-backend-t1/internal"

func NewRepositoryWareHouse() *WarehouseDB {
	return &WarehouseDB{
		DB: make(map[int]*internal.Warehouse),
	}
}

type WarehouseDB struct {
	DB map[int]*internal.Warehouse
}

func (r *WarehouseDB) FindAll() ([]internal.Warehouse, error) {

	return []internal.Warehouse{}, nil
}

func (r *WarehouseDB) FindByID(id int) (internal.Warehouse, error) {

	return internal.Warehouse{}, nil
}

func (r *WarehouseDB) SectionNumberExists(warehouse internal.Warehouse) error {
	return nil
}

func (r *WarehouseDB) Save(warehouse *internal.Warehouse) error {

	return nil
}

func (r *WarehouseDB) Update(warehouse *internal.Warehouse) error {

	return nil
}

func (r *WarehouseDB) Delete(id int) error {

	return nil
}
