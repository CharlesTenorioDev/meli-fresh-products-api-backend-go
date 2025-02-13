package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

// PurchaseOrderJSON is a struct that represents a purchase order in JSON format
type PurchaseOrderJSON struct {
	ID              int    `json:"id"`
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordID int    `json:"product_record_id"`
}

// PurchaseOrderCreateRequest is a struct that represents a purchase order create request
type PurchaseOrderCreateRequest struct {
	OrderNumber     *string `json:"order_number"`
	OrderDate       *string `json:"order_date"`
	TrackingCode    *string `json:"tracking_code"`
	BuyerID         *int    `json:"buyer_id"`
	ProductRecordID *int    `json:"product_record_id"`
}

// NewPurchaseOrderHandler creates a new instance of the purchase order handler
func NewPurchaseOrderHandler(sv internal.PurchaseOrderService) *PurchaseOrderHandler {
	return &PurchaseOrderHandler{
		sv: sv,
	}
}

// PurchaseOrderHandler is the default implementation of the purchase order handler
type PurchaseOrderHandler struct {
	sv internal.PurchaseOrderService
}

// Create creates a new purchase order
// @Summary Create a new purchase order
// @Description Handles the creation of a new purchase order to the database
// @Tags PurchaseOrder
// @Accept json
// @Produce json
// @Param request body handler.PurchaseOrderCreateRequest true "Purchase Order Create Request"
// @Success 201 {object} handler.PurchaseOrderJSON "Created Purchase Order"
// @Failure 400 {object} resterr.RestErr "Invalid data"
// @Failure 422 {object} resterr.RestErr "Purchase Order inputs are Invalid"
// @Failure 404 {object} resterr.RestErr "Product records or Buyer not found"
// @Failure 409 {object} resterr.RestErr "Purchase order number already exists"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/purchase-orders [post]
func (h *PurchaseOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestInput *PurchaseOrderCreateRequest

		// decoding the request
		if err := json.NewDecoder(r.Body).Decode(&requestInput); err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError("Invalid data"))
			return
		}

		// validating the request input required fields
		causes := requestInput.ValidateRequiredFields()
		if len(causes) > 0 {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityWithCausesError(internal.ErrPurchaseOrderUnprocessableEntity.Error(), causes))
			return
		}

		// validating the orderDate field
		orderDate, err := time.Parse(time.DateOnly, *requestInput.OrderDate)
		if err != nil {
			var causes []resterr.Causes
			causes = append(causes, resterr.Causes{
				Field:   "order_date",
				Message: "invalid date format",
			})
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestValidationError(ErrInvalidData, causes))
			return
		}

		// creating the purchase order
		purchaseOrder := &internal.PurchaseOrder{
			ID:              0,
			OrderNumber:     *requestInput.OrderNumber,
			OrderDate:       orderDate,
			TrackingCode:    *requestInput.TrackingCode,
			BuyerID:         *requestInput.BuyerID,
			ProductRecordID: *requestInput.ProductRecordID,
		}

		// saving the purchase order
		if err := h.sv.Save(purchaseOrder); err != nil {
			switch {
			case errors.As(err, &internal.DomainError{}):
				var domainError internal.DomainError
				errors.As(err, &domainError)
				var restCauses []resterr.Causes
				for _, cause := range domainError.Causes {
					restCauses = append(restCauses, resterr.Causes{
						Field:   cause.Field,
						Message: cause.Message,
					})
				}

				response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestValidationError(domainError.Message, restCauses))
			case errors.Is(err, internal.ErrPurchaseOrderConflict):
				response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
			case errors.Is(err, internal.ErrProductRecordsNotFound):
				response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
			case errors.Is(err, service.ErrBuyerNotFound):
				response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
			}

			return
		}

		// creating the response
		purchaseOrderJSON := &PurchaseOrderJSON{
			ID:              purchaseOrder.ID,
			OrderNumber:     purchaseOrder.OrderNumber,
			OrderDate:       purchaseOrder.OrderDate.Format(time.DateOnly),
			TrackingCode:    purchaseOrder.TrackingCode,
			BuyerID:         purchaseOrder.BuyerID,
			ProductRecordID: purchaseOrder.ProductRecordID,
		}

		// sending the response
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": purchaseOrderJSON,
		})
	}
}

// Validating the PurchaseOrderCreateRequest required fields
func (p *PurchaseOrderCreateRequest) ValidateRequiredFields() (causes []resterr.Causes) {
	if p.OrderNumber == nil {
		causes = append(causes, resterr.Causes{
			Field:   "order_number",
			Message: "order number is required",
		})
	}
	if p.OrderDate == nil {
		causes = append(causes, resterr.Causes{
			Field:   "order_date",
			Message: "order date is required",
		})
	}
	if p.TrackingCode == nil {
		causes = append(causes, resterr.Causes{
			Field:   "tracking_code",
			Message: "tracking code is required",
		})
	}
	if p.BuyerID == nil {
		causes = append(causes, resterr.Causes{
			Field:   "buyer_id",
			Message: "buyer id is required",
		})
	}
	if p.ProductRecordID == nil {
		causes = append(causes, resterr.Causes{
			Field:   "product_record_id",
			Message: "product record id is required",
		})
	}
	return
}
