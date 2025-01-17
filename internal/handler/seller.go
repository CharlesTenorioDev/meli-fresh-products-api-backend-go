package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
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

type SellersGetJson struct {
	Id          int    `json:"id"`
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	Locality    int    `json:"locality_id"`
}

type SellersPostJson struct {
	CID         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	Locality    int    `json:"locality_id"`
}

type SellersUpdateJson struct {
	CID         *int    `json:"cid"`
	CompanyName *string `json:"company_name"`
	Address     *string `json:"address"`
	Telephone   *string `json:"telephone"`
	Locality    *int    `json:"locality_id"`
}

// GetAll returns all sellers
func (h *SellerDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := h.sv.FindAll()
		if err != nil {
			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError("sellers not found"))
				return
			}
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		var sellersJson []SellersGetJson
		for i := range all {
			sellersJson = append(sellersJson, SellersGetJson{
				Id:          all[i].ID,
				Cid:         all[i].CID,
				CompanyName: all[i].CompanyName,
				Address:     all[i].Address,
				Telephone:   all[i].Telephone,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": sellersJson,
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
			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
				return
			}
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		var sellerJson = SellersGetJson{
			Id:          seller.ID,
			Cid:         seller.CID,
			CompanyName: seller.CompanyName,
			Address:     seller.Address,
			Telephone:   seller.Telephone,
			Locality:    seller.Locality,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": sellerJson,
		})
	}
}

// Save creates a new seller
func (h *SellerDefault) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var body SellersPostJson
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
			Locality:    body.Locality,
		}

		err = sl.Validate()
		if err != nil {
			restErr := rest_err.NewUnprocessableEntityError(err.Error())
			response.JSON(w, restErr.Code, restErr)
			return
		}

		err = h.sv.Save(sl)
		if err != nil {
			log.Println(err)

			if errors.Is(err, internal.ErrSellerConflict) || errors.Is(err, internal.ErrSellerCIDAlreadyExists) {
				response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
				return
			}

			if errors.Is(err, internal.ErrSellerNotFound) || errors.Is(err, internal.ErrLocalityNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
				return
			}

			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": map[string]any{
				"seller_id": sl.ID,
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

		var body SellersUpdateJson
		err = request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		slPatch := internal.SellerPatch{
			CID:         body.CID,
			CompanyName: body.CompanyName,
			Address:     body.Address,
			Telephone:   body.Telephone,
			Locality:    body.Locality,
		}

		seller, err := h.sv.Update(id, slPatch)
		if err != nil {
			if errors.Is(err, internal.ErrSellerCIDAlreadyExists) {
				response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
				return
			}

			if errors.Is(err, internal.ErrSellerInvalidFields) {
				restErr := rest_err.NewBadRequestError(err.Error())
				response.JSON(w, restErr.Code, restErr)
				return
			}

			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
				return
			}

			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": SellersGetJson{
				Id:          seller.ID,
				Cid:         seller.CID,
				CompanyName: seller.CompanyName,
				Address:     seller.Address,
				Telephone:   seller.Telephone,
				Locality:    seller.Locality,
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
			if errors.Is(err, internal.ErrSellerNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
				return
			}
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}
