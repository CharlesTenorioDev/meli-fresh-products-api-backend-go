package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSellerValidate tests the Validate method of the Seller struct
func TestSellerValidate(t *testing.T) {
	tests := []struct {
		name    string
		seller  Seller
		wantErr bool
		causes  []Causes
	}{
		{
			name: "valid seller",
			seller: Seller{
				ID:          1,
				CID:         1,
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "12 34546-7890",
				Locality:    1,
			},
			wantErr: false,
			causes:  nil,
		},
		{
			name: "invalid seller - negative ID",
			seller: Seller{
				ID:          -1,
				CID:         1,
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "12 34456-7890",
				Locality:    1,
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "id",
					Message: "Seller ID is required",
				},
			},
		},
		{
			name: "invalid seller - missing ID",
			seller: Seller{
				CID:         1,
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "12 34456-7890",
				Locality:    1,
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "id",
					Message: "Seller ID is required",
				},
			},
		},
		{
			name: "invalid seller - missing CID",
			seller: Seller{
				ID:          1,
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "12 34456-7890",
				Locality:    1,
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "cid",
					Message: "Company ID is required",
				},
			},
		},
		{
			name: "invalid seller - missing company name",
			seller: Seller{
				ID:        1,
				CID:       1,
				Address:   "Test Address",
				Telephone: "12 34456-7890",
				Locality:  1,
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "company_name",
					Message: "Company name is required",
				},
			},
		},
		{
			name: "invalid seller - missing address",
			seller: Seller{
				ID:          1,
				CID:         1,
				CompanyName: "Test Company",
				Telephone:   "12 34456-7890",
				Locality:    1,
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "address",
					Message: "Address cannot be empty",
				},
			},
		},
		{
			name: "invalid seller - invalid telephone number",
			seller: Seller{
				ID:          1,
				CID:         1,
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "1234432fds",
				Locality:    1,
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "telephone",
					Message: `Telephone number is invalid, should be formatted as "XX XXXXX-XXXX"`,
				},
			},
		},
		{
			name: "invalid seller - missing locality",
			seller: Seller{
				ID:          1,
				CID:         1,
				CompanyName: "Test Company",
				Address:     "Test Address",
				Telephone:   "21 98888-8888",
			},
			wantErr: true,
			causes: []Causes{
				{
					Field:   "locality_id",
					Message: "Locality ID is required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			causes := tt.seller.Validate()
			if tt.wantErr {
				assert.Equal(t, tt.causes, causes)
			} else {
				assert.Nil(t, causes)
			}
		})
	}
}
