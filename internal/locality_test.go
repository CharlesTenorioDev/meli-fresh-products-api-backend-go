package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocality_Validate(t *testing.T) {
	tests := []struct {
		name     string
		locality internal.Locality
		wantErr  bool
		causes   []internal.Causes
	}{
		{
			name: "valid locality",
			locality: internal.Locality{
				ID:           1,
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			wantErr: false,
			causes:  nil,
		},
		{
			name: "invalid locality - negative ID",
			locality: internal.Locality{
				ID:           -1,
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "locality_id",
					Message: "Locality ID cannot be negative",
				},
			},
		},
		{
			name: "invalid locality - missing ID",
			locality: internal.Locality{
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "locality_id",
					Message: "Locality ID is required",
				},
			},
		},
		{
			name: "invalid locality - missing LocalityName",
			locality: internal.Locality{
				ID:           1,
				ProvinceName: "Test Province",
				CountryName:  "Test Country",
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "locality_name",
					Message: "Locality name is required",
				},
			},
		},
		{
			name: "invalid locality - missing ProvinceName",
			locality: internal.Locality{
				ID:           1,
				LocalityName: "Test Locality",
				CountryName:  "Test Country",
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "province_name",
					Message: "Province name cannot be empty",
				},
			},
		},
		{
			name: "invalid locality - missing CountryName",
			locality: internal.Locality{
				ID:           1,
				LocalityName: "Test Locality",
				ProvinceName: "Test Province",
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "country_name",
					Message: "Country name cannot be empty",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			causes := tt.locality.Validate()
			if tt.wantErr {
				assert.Equal(t, tt.causes, causes)
			} else {
				assert.Nil(t, causes)
			}
		})
	}
}
