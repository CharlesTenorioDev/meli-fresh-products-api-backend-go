package internal

import (
	"errors"
	"time"

	"github.com/meli-fresh-products-api-backend-t1/utils/validator"
)

// PurchaseOrder is a struct that represents a purchase order
type PurchaseOrder struct {
	ID              int
	OrderNumber     string
	OrderDate       time.Time
	TrackingCode    string
	BuyerID         int
	ProductRecordId int
}

var (
	// ErrPurchaseOrderRepositoryNotFound is returned when the purchase order is not found
	ErrPurchaseOrderNotFound = errors.New("purchase order not found")
	// ErrPurchaseOrderRepositoryConflict is returned when the purchase order already exists
	ErrPurchaseOrderConflict = errors.New("purchase order number already exists")
)

// Validate validates the purchase order fields
func (p *PurchaseOrder) Validate() (causes []Causes) {
	if !validator.String(p.OrderNumber, 1, 255) {
		causes = append(causes, Causes{
			Field:   "order_number",
			Message: "Order number is required",
		})
	}
	if !validator.String(p.TrackingCode, 1, 255) {
		causes = append(causes, Causes{
			Field:   "tracking_code",
			Message: "Tracking code is required",
		})
	}
	if validator.IntIsNegative(p.BuyerID) {
		causes = append(causes, Causes{
			Field:   "buyer_id",
			Message: "Buyer ID cannot be negative",
		})
	}
	if validator.IntIsZero(p.BuyerID) {
		causes = append(causes, Causes{
			Field:   "buyer_id",
			Message: "Buyer ID is required",
		})
	}
	if validator.IntIsNegative(p.ProductRecordId) {
		causes = append(causes, Causes{
			Field:   "product_record_id",
			Message: "Product record ID cannot be negative",
		})
	}
	if validator.IntIsZero(p.ProductRecordId) {
		causes = append(causes, Causes{
			Field:   "product_record_id",
			Message: "Product record ID is required",
		})
	}
	return
}

// PurchaseOrdersRepository is an interface that contains the methods that the purchase order repository should support
type PurchaseOrderRepository interface {
	// // FindAll returns all the purchase orders
	// FindAll() ([]PurchaseOrder, error)
	// FindByID returns the purchase order with the given ID
	FindByID(id int) (PurchaseOrder, error)
	// Save saves the given purchase order
	Save(p *PurchaseOrder) error
}

// PurchaseOrderService is an interface that contains the methods that the purchase order service should support
type PurchaseOrderService interface {
	// // FindAll returns all the purchase orders
	// FindAll() ([]PurchaseOrder, error)
	// FindByID returns the purchase order with the given ID
	FindByID(id int) (PurchaseOrder, error)
	// Save saves the given purchase order
	Save(p *PurchaseOrder) error
}
