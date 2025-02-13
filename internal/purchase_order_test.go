package internal_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPurchaseOrder_Validate(t *testing.T) {
	tests := []struct {
		name          string
		purchaseOrder *internal.PurchaseOrder
		wantErr       bool
		causes        []internal.Causes
	}{
		{
			name: "valid purchase order",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(), // Important: Use a valid time.Time
				TrackingCode:    "TRACK-456",
				BuyerID:         1,
				ProductRecordID: 1,
			},
			wantErr: false,
			causes:  nil,
		},
		{
			name: "invalid purchase order - missing order number",
			purchaseOrder: &internal.PurchaseOrder{
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-456",
				BuyerID:         1,
				ProductRecordID: 1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "order_number",
					Message: "order number is required",
				},
			},
		},
		{
			name: "invalid purchase order - order number out of range",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ThisIsAVeryLongOrderNumberThatExceedsTheLimitXXXXXXXXXXXXXXXXXXXX",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-456",
				BuyerID:         1,
				ProductRecordID: 1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "order_number",
					Message: "order number is out of range",
				},
			},
		},
		{
			name: "invalid purchase order - missing tracking code",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(),
				BuyerID:         1,
				ProductRecordID: 1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "tracking_code",
					Message: "tracking code is required",
				},
			},
		},
		{
			name: "invalid purchase order - tracking code out of range",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(),
				TrackingCode:    "ThisIsAVeryLongTrackingCodeThatExceedsTheLimitXXXXXXXXXXXXXXXXXXXXXX",
				BuyerID:         1,
				ProductRecordID: 1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "tracking_code",
					Message: "tracking code is out of range",
				},
			},
		},
		{
			name: "invalid purchase order - zero buyer ID",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-456",
				BuyerID:         0,
				ProductRecordID: 1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "buyer_id",
					Message: "buyer ID is required",
				},
			},
		},
		{
			name: "invalid purchase order - negative buyer ID",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-456",
				BuyerID:         -1,
				ProductRecordID: 1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "buyer_id",
					Message: "buyer ID cannot be negative",
				},
			},
		},
		{
			name: "invalid purchase order - zero product record ID",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-456",
				BuyerID:         1,
				ProductRecordID: 0,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "product_record_id",
					Message: "product record ID is required",
				},
			},
		},
		{
			name: "invalid purchase order - negative product record ID",
			purchaseOrder: &internal.PurchaseOrder{
				OrderNumber:     "ORDER-123",
				OrderDate:       time.Now(),
				TrackingCode:    "TRACK-456",
				BuyerID:         1,
				ProductRecordID: -1,
			},
			wantErr: true,
			causes: []internal.Causes{
				{
					Field:   "product_record_id",
					Message: "product record ID cannot be negative",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			causes := tt.purchaseOrder.Validate()
			if tt.wantErr {
				assert.Equal(t, tt.causes, causes)
			} else {
				assert.Nil(t, causes)
			}
		})
	}
}
