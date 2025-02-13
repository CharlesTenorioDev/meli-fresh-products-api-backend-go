package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCarries_Ok(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) *internal.Carries
		expectedOutput bool
	}{
		{
			name: "Should return false when Cid is empty",
			setup: func(t *testing.T) *internal.Carries {
				return &internal.Carries{
					Cid:         "",
					CompanyName: "Test Company",
					Address:     "Test Address",
					PhoneNumber: "123-456-7890",
					LocalityID:  1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when CompanyName is empty",
			setup: func(t *testing.T) *internal.Carries {
				return &internal.Carries{
					Cid:         "Test Cid",
					CompanyName: "",
					Address:     "Test Address",
					PhoneNumber: "123-456-7890",
					LocalityID:  1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when Address is empty",
			setup: func(t *testing.T) *internal.Carries {
				return &internal.Carries{
					Cid:         "Test Cid",
					CompanyName: "Test Company",
					Address:     "",
					PhoneNumber: "123-456-7890",
					LocalityID:  1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when PhoneNumber is empty",
			setup: func(t *testing.T) *internal.Carries {
				return &internal.Carries{
					Cid:         "Test Cid",
					CompanyName: "Test Company",
					Address:     "Test Address",
					PhoneNumber: "",
					LocalityID:  1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return false when LocalityID is negative",
			setup: func(t *testing.T) *internal.Carries {
				return &internal.Carries{
					Cid:         "Test Cid",
					CompanyName: "Test Company",
					Address:     "Test Address",
					PhoneNumber: "123-456-7890",
					LocalityID:  -1,
				}
			},
			expectedOutput: false,
		},
		{
			name: "Should return true when all fields are valid",
			setup: func(t *testing.T) *internal.Carries {
				return &internal.Carries{
					Cid:         "Test Cid",
					CompanyName: "Test Company",
					Address:     "Test Address",
					PhoneNumber: "123-456-7890",
					LocalityID:  1,
				}
			},
			expectedOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.setup(t)

			actualOutput := c.Ok()

			assert.Equal(t, tt.expectedOutput, actualOutput)
		})
	}
}
