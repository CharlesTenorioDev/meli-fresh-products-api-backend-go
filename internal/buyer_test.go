package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuyer_Parse(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) *internal.Buyer
		expectedOutput bool
	}{
		{
			name: "Should return true when valid input",
			setup: func(t *testing.T) *internal.Buyer {
				return &internal.Buyer{
					ID:           1,
					CardNumberID: "1212312321",
					FirstName:    "Cesar",
					LastName:     "C C R Pontes",
				}
			},
			expectedOutput: true,
		},
		{
			name: "Should return false when invalid CardNumberId input",
			setup: func(t *testing.T) *internal.Buyer {
				return &internal.Buyer{
					ID:           1,
					CardNumberID: "",
					FirstName:    "Cesar",
					LastName:     "C C R Pontes",
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when invalid FirstName input",
			setup: func(t *testing.T) *internal.Buyer {
				return &internal.Buyer{
					ID:           1,
					CardNumberID: "1212312321",
					FirstName:    "",
					LastName:     "C C R Pontes",
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when invalid LastName input",
			setup: func(t *testing.T) *internal.Buyer {
				return &internal.Buyer{
					ID:           1,
					CardNumberID: "1212312321",
					FirstName:    "Cesar",
					LastName:     "",
				}
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.setup(t)
			actualOutput := b.Parse()
			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}

func TestBuyerPatch_Patch(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) (internal.BuyerPatch, *internal.Buyer)
		expectedOutput *internal.Buyer
	}{
		{
			name: "Should update all fields",
			setup: func(t *testing.T) (internal.BuyerPatch, *internal.Buyer) {
				return internal.BuyerPatch{
						CardNumberID: strPtr("10212121"),
						FirstName:    strPtr("Jack"),
						LastName:     strPtr("C C R Salmon"),
					}, &internal.Buyer{
						ID:           1,
						CardNumberID: "124324432",
						FirstName:    "Jackie",
						LastName:     "C C R Pontes",
					}
			},
			expectedOutput: &internal.Buyer{
				ID:           1,
				CardNumberID: "10212121",
				FirstName:    "Jack",
				LastName:     "C C R Salmon",
			},
		},
		{
			name: "Should update only CardNumberId",
			setup: func(t *testing.T) (internal.BuyerPatch, *internal.Buyer) {
				return internal.BuyerPatch{
						CardNumberID: strPtr("10212121"),
					}, &internal.Buyer{
						ID:           1,
						CardNumberID: "124324432",
						FirstName:    "Jackie",
						LastName:     "C C R Pontes",
					}
			},
			expectedOutput: &internal.Buyer{
				ID:           1,
				CardNumberID: "10212121",
				FirstName:    "Jackie",
				LastName:     "C C R Pontes",
			},
		},
		{
			name: "Should update only FirstName",
			setup: func(t *testing.T) (internal.BuyerPatch, *internal.Buyer) {
				return internal.BuyerPatch{
						FirstName: strPtr("Jack"),
					}, &internal.Buyer{
						ID:           1,
						CardNumberID: "124324432",
						FirstName:    "Jackie",
						LastName:     "C C R Pontes",
					}
			},
			expectedOutput: &internal.Buyer{
				ID:           1,
				CardNumberID: "124324432",
				FirstName:    "Jack",
				LastName:     "C C R Pontes",
			},
		},
		{
			name: "Should update only LastName",
			setup: func(t *testing.T) (internal.BuyerPatch, *internal.Buyer) {
				return internal.BuyerPatch{
						LastName: strPtr("C C R Salmon"),
					}, &internal.Buyer{
						ID:           1,
						CardNumberID: "124324432",
						FirstName:    "Jackie",
						LastName:     "C C R Pontes",
					}
			},
			expectedOutput: &internal.Buyer{
				ID:           1,
				CardNumberID: "124324432",
				FirstName:    "Jackie",
				LastName:     "C C R Salmon",
			},
		},
		{
			name: "Shouldn't update any field",
			setup: func(t *testing.T) (internal.BuyerPatch, *internal.Buyer) {
				return internal.BuyerPatch{
						CardNumberID: nil,
						FirstName:    nil,
						LastName:     nil,
					}, &internal.Buyer{
						ID:           1,
						CardNumberID: "124324432",
						FirstName:    "Jackie",
						LastName:     "C C R Pontes",
					}
			},
			expectedOutput: &internal.Buyer{
				ID:           1,
				CardNumberID: "124324432",
				FirstName:    "Jackie",
				LastName:     "C C R Pontes",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp, b := tt.setup(t)

			bp.Patch(b)

			assert.Equal(t, tt.expectedOutput, b)
		})
	}

}
