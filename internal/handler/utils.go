package handler

import (
	"encoding/json"
	"net/http"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type RequestBodySection struct {
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}

type Data struct {
	ID                 int     `json:"id"`
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}

type ResponseBodySection struct {
	Data *Data `json:"data,omitempty"`
}

type ResponseBodySections struct {
	Data *[]Data `json:"data,omitempty"`
}

type ResponseBodyError struct {
	Data string `json:"data"`
}

func HandleError(w http.ResponseWriter, message string, statusCode int) {
	body := &ResponseBodyError{
		Data: message,
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func RespondWithSection(w http.ResponseWriter, section internal.Section, statusCode int) {
	dt := Data{
		ID:                 section.ID,
		SectionNumber:      section.SectionNumber,
		CurrentTemperature: section.CurrentTemperature,
		MinimumTemperature: section.MinimumTemperature,
		CurrentCapacity:    section.CurrentCapacity,
		MinimumCapacity:    section.MinimumCapacity,
		MaximumCapacity:    section.MaximumCapacity,
		WarehouseID:        section.WarehouseID,
		ProductTypeID:      section.ProductTypeID,
	}

	body := ResponseBodySection{
		Data: &dt,
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func RespondWithSections(w http.ResponseWriter, sections []internal.Section, statusCode int) {
	var datas []Data
	for _, value := range sections {
		dt := Data{
			ID:                 value.ID,
			SectionNumber:      value.SectionNumber,
			CurrentTemperature: value.CurrentTemperature,
			MinimumTemperature: value.MinimumTemperature,
			CurrentCapacity:    value.CurrentCapacity,
			MinimumCapacity:    value.MinimumCapacity,
			MaximumCapacity:    value.MaximumCapacity,
			WarehouseID:        value.WarehouseID,
			ProductTypeID:      value.ProductTypeID,
		}

		datas = append(datas, dt)
	}

	body := ResponseBodySections{
		Data: &datas,
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
