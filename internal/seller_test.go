package internal_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeller_Validate(t *testing.T) {
	tests := []struct {
		name    string
		seller  internal.Seller
		wantErr bool
		err     error
	}{
		{
			name: "valid seller",
			seller: internal.Seller{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
				Locality:    1,
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "invalid seller - missing CID",
			seller: internal.Seller{
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
				Locality:    1,
			},
			wantErr: true,
			err:     errors.Join(errors.New("seller.CID is required")),
		},
		{
			name: "invalid seller - missing CompanyName",
			seller: internal.Seller{
				CID:       123,
				Address:   "Rua 1",
				Telephone: "1234567890",
				Locality:  1,
			},
			wantErr: true,
			err:     errors.Join(errors.New("seller.CompanyName is required")),
		},
		{
			name: "invalid seller - missing Address",
			seller: internal.Seller{
				CID:         123,
				CompanyName: "Test Seller",
				Telephone:   "1234567890",
				Locality:    1,
			},
			wantErr: true,
			err:     errors.Join(errors.New("seller.Address is required")),
		},
		{
			name: "invalid seller - missing Telephone",
			seller: internal.Seller{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Locality:    1,
			},
			wantErr: true,
			err:     errors.Join(errors.New("seller.Telephone is required")),
		},
		{
			name: "invalid seller - missing Locality",
			seller: internal.Seller{
				CID:         123,
				CompanyName: "Test Seller",
				Address:     "Rua 1",
				Telephone:   "1234567890",
			},
			wantErr: true,
			err:     errors.Join(errors.New("seller.Locality is required")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.seller.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
