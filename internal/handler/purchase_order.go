package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
)

// PurchaseOrderJSON is a struct that represents a purchase order in JSON format
type PurchaseOrderJSON struct {
	ID              int    `json:"id"`
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordId int    `json:"product_record_id"`
}

// PurchaseOrderCreateRequest is a struct that represents a purchase order create request
type PurchaseOrderCreateRequest struct {
	OrderNumber     string `json:"order_number"`
	OrderDate       string `json:"order_date"`
	TrackingCode    string `json:"tracking_code"`
	BuyerID         int    `json:"buyer_id"`
	ProductRecordId int    `json:"product_record_id"`
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
func (h *PurchaseOrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestInput *PurchaseOrderCreateRequest

		// decoding the request
		if err := json.NewDecoder(r.Body).Decode(&requestInput); err != nil {
			response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
			return
		}

		// validating the orderDate field
		orderDate, err := time.Parse(time.DateOnly, requestInput.OrderDate)
		if err != nil {
			var causes []rest_err.Causes

			causes = append(causes, rest_err.Causes{
				Field:   "order_date",
				Message: "invalid date format",
			})

			response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestValidationError(ErrInvalidData, causes))
			return
		}

		// creating the purchase order
		purchaseOrder := &internal.PurchaseOrder{
			ID:              0,
			OrderNumber:     requestInput.OrderNumber,
			OrderDate:       orderDate,
			TrackingCode:    requestInput.TrackingCode,
			BuyerID:         requestInput.BuyerID,
			ProductRecordId: requestInput.ProductRecordId,
		}

		// saving the purchase order
		if err := h.sv.Save(purchaseOrder); err != nil {
			switch {
			case errors.As(err, &internal.DomainError{}):
				var domainError internal.DomainError
				errors.As(err, &domainError)

				var restCauses []rest_err.Causes
				for _, cause := range domainError.Causes {
					restCauses = append(restCauses, rest_err.Causes{
						Field:   cause.Field,
						Message: cause.Message,
					})
				}
				response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewBadRequestValidationError(domainError.Message, restCauses))
			case errors.Is(err, internal.ErrPurchaseOrderConflict):
				response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
			case errors.Is(err, service.BuyerNotFound):
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(ErrInternalServer))
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
			ProductRecordId: purchaseOrder.ProductRecordId,
		}

		// sending the response
		response.JSON(w, http.StatusCreated, map[string]any{
			"data": purchaseOrderJSON,
		})
	}
}
