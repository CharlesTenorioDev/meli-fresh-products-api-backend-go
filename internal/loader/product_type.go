package loader

import (
	"encoding/json"
	"io"
	"os"
)

const (
	localFileProdTypeJSON = "db/product_type.json"
)

type ProductType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ReadAllProductsTypeToFile() ([]*ProductType, error) {
	var productTypeList []*ProductType

	file, err := os.Open(localFileProdTypeJSON)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(localFileProdTypeJSON)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			initialData := []ProductType{}

			writer := json.NewEncoder(file)
			if err := writer.Encode(initialData); err != nil {
				return nil, err
			}

			return productTypeList, nil
		}

		return nil, err
	}
	defer file.Close()

	reader := json.NewDecoder(file)

	err = reader.Decode(&productTypeList)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}

		return nil, err
	}

	return productTypeList, nil
}
