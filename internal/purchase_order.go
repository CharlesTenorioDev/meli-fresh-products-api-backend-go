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
	ProductRecordID int
}

var (
	// ErrPurchaseOrderRepositoryNotFound is returned when the purchase order is not found
	ErrPurchaseOrderNotFound = errors.New("purchase order not found")
	// ErrPurchaseOrderRepositoryConflict is returned when the purchase order already exists
	ErrPurchaseOrderConflict = errors.New("purchase order number already exists")
	// ErrPurchaseOrderUnprocessableEntity is returned when the purchase order is unprocessable
	ErrPurchaseOrderUnprocessableEntity = errors.New("purchase order inputs are missing")
	// ErrPurchaseOrderBadRequest is returned when the purchase order request is bad
	ErrPurchaseOrderBadRequest = errors.New("purchase order inputs are invalid")
)

// Validate validates the business rules of the purchase order
func (p *PurchaseOrder) Validate() (causes []Causes) {
	if validator.BlankString(p.OrderNumber) {
		causes = append(causes, Causes{
			Field:   "order_number",
			Message: "order number is required",
		})
	}

	if !validator.String(p.OrderNumber, 1, 50) && !validator.BlankString(p.OrderNumber) {
		causes = append(causes, Causes{
			Field:   "order_number",
			Message: "order number is out of range",
		})
	}

	if validator.BlankString(p.TrackingCode) {
		causes = append(causes, Causes{
			Field:   "tracking_code",
			Message: "tracking code is required",
		})
	}

	if !validator.String(p.TrackingCode, 1, 50) && !validator.BlankString(p.TrackingCode) {
		causes = append(causes, Causes{
			Field:   "tracking_code",
			Message: "tracking code is out of range",
		})
	}

	if validator.IntIsZero(p.BuyerID) {
		causes = append(causes, Causes{
			Field:   "buyer_id",
			Message: "buyer ID is required",
		})
	}

	if validator.IntIsNegative(p.BuyerID) {
		causes = append(causes, Causes{
			Field:   "buyer_id",
			Message: "buyer ID cannot be negative",
		})
	}

	if validator.IntIsZero(p.ProductRecordID) {
		causes = append(causes, Causes{
			Field:   "product_record_id",
			Message: "product record ID is required",
		})
	}

	if validator.IntIsNegative(p.ProductRecordID) {
		causes = append(causes, Causes{
			Field:   "product_record_id",
			Message: "product record ID cannot be negative",
		})
	}

	return causes
}

// PurchaseOrderRepository is an interface that contains the methods that the purchase order repository should support
type PurchaseOrderRepository interface {
	// FindByID returns the purchase order with the given ID
	FindByID(id int) (PurchaseOrder, error)
	// Save saves the given purchase order
	Save(p *PurchaseOrder) error
}

// PurchaseOrderService is an interface that contains the methods that the purchase order service should support
type PurchaseOrderService interface {
	// FindByID returns the purchase order with the given ID
	FindByID(id int) (PurchaseOrder, error)
	// Save saves the given purchase order
	Save(p *PurchaseOrder) error
}
