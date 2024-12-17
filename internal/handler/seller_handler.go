package handler

import (
	"errors"
	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/dto"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
	"net/http"
	"strconv"
)

// NewSellerDefault creates a new instance of the seller handler
func NewSellerDefault(sv internal.SellerService) *SellerDefault {
	return &SellerDefault{
		sv: sv,
	}
}

// SellerDefault is the default implementation of the seller handler
type SellerDefault struct {
	// sv is the service used by the handler
	sv internal.SellerService
}

// GetAll returns all sellers
func (h *SellerDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)

			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError("sellers not found"))
			}

			return
		}

		var sellersJson []dto.SellersGetDto
		for i := range all {
			sellersJson = append(sellersJson, dto.SellersGetDto{
				Id:          all[i].ID,
				Cid:         all[i].CID,
				CompanyName: all[i].CompanyName,
				Address:     all[i].Address,
				Telephone:   all[i].Telephone,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    sellersJson,
		})
	}
}

// GetByID returns a seller
func (h *SellerDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		seller, err := h.sv.FindByID(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)

			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			}

			return
		}

		var sellerJson = dto.SellersGetDto{
			Id:          seller.ID,
			Cid:         seller.CID,
			CompanyName: seller.CompanyName,
			Address:     seller.Address,
			Telephone:   seller.Telephone,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    sellerJson,
		})
	}
}

// Save creates a new seller
func (h *SellerDefault) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body dto.SellersPostDto
		err := request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		sl := &internal.Seller{
			CID:         body.CID,
			CompanyName: body.CompanyName,
			Address:     body.Address,
			Telephone:   body.Telephone,
		}

		err = sl.Validate()
		if err != nil {
			restErr := rest_err.NewUnprocessableEntityError(err.Error())
			response.JSON(w, restErr.Code, restErr)
			return
		}

		id, err := h.sv.Save(sl)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)

			if errors.Is(err, internal.ErrSellerConflict) || errors.Is(err, internal.ErrSellerCIDAlreadyExists) {
				response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
			}

			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			}

			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data": map[string]any{
				"seller_id": id,
			},
		})
	}
}

// Update updates a seller
func (h *SellerDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		actualSeller, err := h.sv.FindByID(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)

			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			}

			return
		}

		var body dto.SellersUpdateDto
		err = request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
		}

		if body.CompanyName != nil {
			actualSeller.CompanyName = *body.CompanyName
		}

		if body.Address != nil {
			actualSeller.Address = *body.Address
		}

		if body.Telephone != nil {
			actualSeller.Telephone = *body.Telephone
		}

		err = actualSeller.Validate()
		if err != nil {
			restErr := rest_err.NewUnprocessableEntityError(err.Error())
			response.JSON(w, restErr.Code, restErr)
			return
		}

		err = h.sv.Update(&actualSeller)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)

			if errors.Is(err, internal.ErrSellerCIDAlreadyExists) {
				response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
			}

			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			}

			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data": dto.SellersGetDto{
				Id:          actualSeller.ID,
				Cid:         actualSeller.CID,
				CompanyName: actualSeller.CompanyName,
				Address:     actualSeller.Address,
				Telephone:   actualSeller.Telephone,
			},
		})

	}
}

// Delete deletes a seller
func (h *SellerDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		err = h.sv.Delete(id)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			}
			return
		}

		response.JSON(w, http.StatusNoContent, map[string]any{
			"message": "success",
			"data":    nil,
		})
	}
}
