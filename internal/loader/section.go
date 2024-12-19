package loader

import (
	"encoding/json"
	"io"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

const (
	localFileJson = "db/section.json"
)

func ReadAllSectionsToFile() ([]*internal.Section, error) {
	var sectionList []*internal.Section

	file, err := os.Open(localFileJson)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(localFileJson)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			initialData := []internal.Section{}
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
