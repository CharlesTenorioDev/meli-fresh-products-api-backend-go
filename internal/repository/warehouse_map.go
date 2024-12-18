package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type WarehouseAttributesMap struct {
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
}

type RepositoryWarehouseMap struct {
	db     map[int]WarehouseAttributesMap
	lastID int
}

// RepositoryWarehouseMap Builder that reads a JSON file and returns a map of warehouses
func NewRepositoryWarehouse(db map[int]WarehouseAttributesMap, filePath string) *RepositoryWarehouseMap {
	// If the db is not nil, return the map
	if db != nil {
		return &RepositoryWarehouseMap{db: db, lastID: len(db)}
	}

	// If the db is nil, read the JSON file and return the map
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Decoding JSON file
	var fileData []WarehouseAttributesMap
	err = json.NewDecoder(file).Decode(&fileData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Serialization
	db = make(map[int]WarehouseAttributesMap)
	for key, value := range fileData {
		db[key+1] = WarehouseAttributesMap{
			WarehouseCode:      value.WarehouseCode,
			Address:            value.Address,
			Telephone:          value.Telephone,
			MinimumCapacity:    value.MinimumCapacity,
			MinimumTemperature: value.MinimumTemperature,
		}
	}
	return &RepositoryWarehouseMap{db: db, lastID: len(db)}
}

// FindAll is a method that returns all Warehouses
func (r *RepositoryWarehouseMap) FindAll() (w []internal.Warehouse, err error) {
	w = make([]internal.Warehouse, 0, len(r.db))

	for key, value := range r.db {
		w = append(w, internal.Warehouse{
			ID:                 key,
			WarehouseCode:      value.WarehouseCode,
			Address:            value.Address,
			Telephone:          value.Telephone,
			MinimumCapacity:    value.MinimumCapacity,
			MinimumTemperature: value.MinimumTemperature,
		})
	}

	// Sorting all warehouses by ID
	for i := 0; i < len(w); i++ {
		for j := i + 1; j < len(w); j++ {
			if w[i].ID > w[j].ID {
				w[i], w[j] = w[j], w[i]
			}
		}
	}
	return w, nil
}

// FindByID is a method that returns a Warehouses by its ID
func (r *RepositoryWarehouseMap) FindByID(id int) (w internal.Warehouse, err error) {
	wh, ok := r.db[id]
	if !ok {
		err = fmt.Errorf("%w: %d", internal.ErrWarehouseRepositoryNotFound, id)
		return
	}

	// returning a domain object
	w = internal.Warehouse{
		ID:                 id,
		WarehouseCode:      wh.WarehouseCode,
		Address:            wh.Address,
		Telephone:          wh.Telephone,
		MinimumCapacity:    wh.MinimumCapacity,
		MinimumTemperature: wh.MinimumTemperature,
	}

	return
}

// Save is a method that saves a Warehouse
func (r *RepositoryWarehouseMap) Save(warehouse *internal.Warehouse) (err error) {
	wh := WarehouseAttributesMap{
		WarehouseCode:      warehouse.WarehouseCode,
		Address:            warehouse.Address,
		Telephone:          warehouse.Telephone,
		MinimumCapacity:    warehouse.MinimumCapacity,
		MinimumTemperature: warehouse.MinimumTemperature,
	}

	// increment the lastID
	r.lastID++
	// save the warehouse
	r.db[r.lastID] = wh
	// set the id of the warehouse
	(*warehouse).ID = r.lastID

	return nil
}

// Update is a method that updates a Warehouse
func (r *RepositoryWarehouseMap) Update(warehouse *internal.Warehouse) (err error) {

	wh := WarehouseAttributesMap{
		WarehouseCode:      warehouse.WarehouseCode,
		Address:            warehouse.Address,
		Telephone:          warehouse.Telephone,
		MinimumCapacity:    warehouse.MinimumCapacity,
		MinimumTemperature: warehouse.MinimumTemperature,
	}

	// update the warehouse
	r.db[warehouse.ID] = wh

	return nil
}

// Delete is a method that deletes a Warehouse
func (r *RepositoryWarehouseMap) Delete(id int) (err error) {
	// Deleting the warehouse
	delete(r.db, id)
	return
}
