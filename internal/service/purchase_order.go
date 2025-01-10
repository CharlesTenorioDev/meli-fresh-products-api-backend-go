package service

import (
	"fmt"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

// NewPurchaseOrderService creates a new instance of the purchase order service
func NewPurchaseOrderService(rpPurchaseOrder internal.PurchaseOrderRepository, rpProductRecords internal.ProductRecordsRepository, svBuyer internal.BuyerService) *PurchaseOrderService {
	return &PurchaseOrderService{
		rpPurchaseOrder:  rpPurchaseOrder,
		rpProductRecords: rpProductRecords,
		svBuyer:          svBuyer,
	}
}

// PurchaseOrderService is the implementation of the purchase order service
type PurchaseOrderService struct {
	rpPurchaseOrder  internal.PurchaseOrderRepository
	rpProductRecords internal.ProductRecordsRepository
	svBuyer          internal.BuyerService
}

// FindById returns a purchase order
func (s *PurchaseOrderService) FindByID(id int) (p internal.PurchaseOrder, err error) {
	p, err = s.rpPurchaseOrder.FindByID(id)
	return
}

// Save creates a new purchase order
func (s *PurchaseOrderService) Save(p *internal.PurchaseOrder) (err error) {
	// Validate the purchase order entity
	causes := p.Validate()
	if len(causes) > 0 {
		return internal.DomainError{
			Message: "Purchase Order inputs are Invalid",
			Causes:  causes,
		}
	}

	// Check if the product records exists
	_, err = s.rpProductRecords.FindByID(p.ProductRecordId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Check if the buyer exists
	_, err = s.svBuyer.FindByID(p.BuyerID)
	if err != nil {
		return err
	}

	// Save the purchase order
	err = s.rpPurchaseOrder.Save(p)
	return
}
