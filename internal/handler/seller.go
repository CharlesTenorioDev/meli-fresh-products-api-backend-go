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
// @Summary Retrieve all sellers
// @Description Fetches a list of all sellers in the database
// @Tags Seller
// @Produce json
// @Success 200 {object} []SellersGetJson "List of sellers"
// @Failure 404 {object} rest_err.RestErr "Sellers not found"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sellers [get]
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

// GetByID returns a seller by ID
// @Summary Retrieve a seller by ID
// @Description Fetches a seller's details based on the provided ID
// @Tags Seller
// @Produce json
// @Param id path int true "Seller ID"
// @Success 200 {object} SellersGetJson "Seller data"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Seller Not Found"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sellers/{id} [get]
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
// @Summary Create a new seller
// @Description Adds a new seller to the system with the provided details on the request body.
// @Tags Seller
// @Accept json
// @Produce json
// @Param seller body SellersPostJson true "Seller Create Request"
// @Success 201 {object} map[string]interface{} "Created Seller Id"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Seller not found" or "Locality not found"
// @Failure 409 {object} rest_err.RestErr "Seller already exists" or "Seller with this CID already exists"
// @Failure 422 {object} rest_err.RestErr "Unprocessable Entity"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sellers [post]
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

// Update updates a seller's details
// @Summary Update seller details
// @Description Modify the information of an existing seller
// @Tags Seller
// @Accept json
// @Produce json
// @Param id path int true "Seller ID"
// @Param seller body SellersUpdateJson true "Seller Update Request"
// @Success 200 {object} SellersGetJson "Updated Seller data"
// @Failure 400 {object} rest_err.RestErr "Seller invalid fields"
// @Failure 404 {object} rest_err.RestErr "Seller not found"
// @Failure 409 {object} rest_err.RestErr "Seller with this CID already exists"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sellers/{id} [patch]
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
// @Summary Delete a seller
// @Description Removes a seller from the system based on the provided Id
// @Tags Seller
// @Produce json
// @Param id path int true "Seller ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Seller not found"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sellers/{id} [delete]
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
