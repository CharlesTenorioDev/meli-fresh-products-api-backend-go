package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWarehouse_Validate(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) *internal.Warehouse
		expectedOutput []internal.Causes
	}{
		{
			name: "",
			setup: func(t *testing.T) *internal.Warehouse {
				return &internal.Warehouse{
					ID:                 0,
					WarehouseCode:      "",
					Address:            "",
					Telephone:          "",
					MinimumCapacity:    0,
					MinimumTemperature: 0,
				}
			},
			expectedOutput: []internal.Causes{
				{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ware := tt.setup(t)

			actualOutput := ware.Validate()

			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}
