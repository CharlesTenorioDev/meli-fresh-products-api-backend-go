package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWarehouseValidate(t *testing.T) {
	tests := []struct {
		name      string
		warehouse internal.Warehouse
		wantErr   bool
		causes    []internal.Causes
	}{
		{
			name: "valid warehouse",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: false,
			causes:  nil,
		},
		{
			name: "invalid warehouse - missing warehouse code",
			warehouse: internal.Warehouse{
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "warehouse_code",
					Message: "warehouse code is required",
				},
			},
		},
		{
			name: "invalid warehouse - warehouse code out of range",
			warehouse: internal.Warehouse{
				WarehouseCode:      "ThisIsAVeryLongWarehouseCodeThatExceedsTheLimitXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "warehouse_code",
					Message: "warehouse code is out of range",
				},
			},
		},
		{
			name: "invalid warehouse - missing address",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "address",
					Message: "address is required",
				},
			},
		},
		{
			name: "invalid warehouse - address out of range",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "ThisIsAVeryLongAddressThatExceedsTheLimitXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "address",
					Message: "address is out of range",
				},
			},
		},

		{
			name: "invalid warehouse - missing telephone",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "telephone",
					Message: "telephone is required",
				},
				{
					Field:   "telephone",
					Message: "telephone number is invalid, should be formatted as XX XXXXX-XXXX",
				},
			},
		},
		{
			name: "invalid warehouse - invalid telephone",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12345678901", // Invalid format
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "telephone",
					Message: `telephone number is invalid, should be formatted as XX XXXXX-XXXX`,
				},
			},
		},
		{
			name: "invalid warehouse - telephone out of range",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12 34567-890111", // Invalid format
				MinimumCapacity:    100,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "telephone",
					Message: "telephone number is invalid, should be formatted as XX XXXXX-XXXX",
				},
			},
		},
		{
			name: "invalid warehouse - zero minimum capacity",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    0,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "minimum_capacity",
					Message: "minimum capacity is required",
				},
			},
		},
		{
			name: "invalid warehouse - negative minimum capacity",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    -10,
				MinimumTemperature: 20.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "minimum_capacity",
					Message: "minimum capacity cannot be negative",
				},
			},
		},
		{
			name: "invalid warehouse - minimum temperature out of range (too low)",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: -300.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "minimum_temperature",
					Message: "minimum temperature is out of range",
				},
			},
		},
		{
			name: "invalid warehouse - minimum temperature out of range (too high)",
			warehouse: internal.Warehouse{
				WarehouseCode:      "WH-123",
				Address:            "Test Address",
				Telephone:          "12 34567-8901",
				MinimumCapacity:    100,
				MinimumTemperature: 1200.0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "minimum_temperature",
					Message: "minimum temperature is out of range",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			causes := tt.warehouse.Validate()
			if tt.wantErr {
				assert.Equal(t, tt.causes, causes)
			} else {
				assert.Nil(t, causes)
			}
		})
	}
}
