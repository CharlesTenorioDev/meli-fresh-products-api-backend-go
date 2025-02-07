package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSection_Ok(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) *internal.Section
		expectedOutput bool
	}{
		{
			name: "Should return false when SectionNumber is zero",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      0,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when CurrentTemperature is below absolute zero",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: -274,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when MinimumTemperature is below absolute zero",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: -274,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when CurrentCapacity is negative",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    -1,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when MinimumCapacity is negative",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    -1,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when MaximumCapacity is negative",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    -1,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when WarehouseID is zero",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        0,
					ProductTypeID:      1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when ProductTypeID is zero",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      0,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return true when all fields are valid",
			setup: func(t *testing.T) *internal.Section {
				return &internal.Section{
					SectionNumber:      1,
					CurrentTemperature: 25,
					MinimumTemperature: 20,
					CurrentCapacity:    100,
					MinimumCapacity:    50,
					MaximumCapacity:    200,
					WarehouseID:        1,
					ProductTypeID:      1,
				}
			},
			expectedOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.setup(t)

			actualOutput := s.Ok()

			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}
