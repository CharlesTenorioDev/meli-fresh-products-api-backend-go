package loader

import (
	"encoding/json"
	"io"
	"os"
)

const (
	localFileJson = "../db/section.json"
)

type SectionJson struct {
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

func ReadAllSectionsToFile() ([]*SectionJson, error) {
	var sectionList []*SectionJson

	file, err := os.Open(localFileJson)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(localFileJson)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			initialData := []SectionJson{}
			writer := json.NewEncoder(file)
			if err := writer.Encode(initialData); err != nil {
				return nil, err
			}

			return sectionList, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := json.NewDecoder(file)
	err = reader.Decode(&sectionList)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	return sectionList, nil
}
