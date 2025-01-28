package loader

import (
	"encoding/json"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func ReadAllSectionsToFile(dbPath string) ([]*internal.Section, error) {
	var sectionList []*internal.Section

	file, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&sectionList)
	if err != nil {
		return nil, err
	}

	return sectionList, nil
}
