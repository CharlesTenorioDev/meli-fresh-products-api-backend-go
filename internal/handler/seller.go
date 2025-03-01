package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/request"
	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
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

type SellersGetJSON struct {
	ID          int    `json:"id"`
	CID         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	Locality    int    `json:"locality_id"`
}

type SellersPostJSON struct {
	CID         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	Locality    int    `json:"locality_id"`
}

type SellersUpdateJSON struct {
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
// @Success 200 {object} []SellersGetJSON "List of sellers"
// @Failure 404 {object} resterr.RestErr "Sellers not found"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/sellers [get]
func (h *SellerDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := h.sv.FindAll()
		if err != nil {
			h.handleError(w, err)
			return
		}

		var sellerJSON []SellersGetJSON
		for i := range all {
			sellerJSON = append(sellerJSON, SellersGetJSON{
				ID:          all[i].ID,
				CID:         all[i].CID,
				CompanyName: all[i].CompanyName,
				Address:     all[i].Address,
				Telephone:   all[i].Telephone,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": sellerJSON,
		})
	}
}

// GetByID returns a seller by ID
// @Summary Retrieve a seller by ID
// @Description Fetches a seller's details based on the provided ID
// @Tags Seller
// @Produce json
// @Param id path int true "Seller ID"
// @Success 200 {object} SellersGetJSON "Seller data"
// @Failure 400 {object} resterr.RestErr "Bad Request"
// @Failure 404 {object} resterr.RestErr "Seller Not Found"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/sellers/{id} [get]
func (h *SellerDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
			return
		}

		seller, err := h.sv.FindByID(id)
		if err != nil {
			h.handleError(w, err)
			return
		}

		var sellerJSON = SellersGetJSON{
			ID:          seller.ID,
			CID:         seller.CID,
			CompanyName: seller.CompanyName,
			Address:     seller.Address,
			Telephone:   seller.Telephone,
			Locality:    seller.Locality,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": sellerJSON,
		})
	}
}

// Save creates a new seller
// @Summary Create a new seller
// @Description Adds a new seller to the system with the provided details on the request body.
// @Tags Seller
// @Accept json
// @Produce json
// @Param seller body SellersPostJSON true "Seller Create Request"
// @Success 201 {object} map[string]interface{} "Created Seller Id"
// @Failure 400 {object} resterr.RestErr "Bad Request"
// @Failure 404 {object} resterr.RestErr "Seller not found" or "Locality not found"
// @Failure 409 {object} resterr.RestErr "Seller already exists" or "Seller with this CID already exists"
// @Failure 422 {object} resterr.RestErr "Unprocessable Entity"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/sellers [post]
func (h *SellerDefault) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body SellersPostJSON

		err := request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError("request json invalid"))
			return
		}

		sl := &internal.Seller{
			CID:         body.CID,
			CompanyName: body.CompanyName,
			Address:     body.Address,
			Telephone:   body.Telephone,
			Locality:    body.Locality,
		}

		err = h.sv.Save(sl)
		if err != nil {
			h.handleError(w, err)
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
// @Param seller body SellersUpdateJSON true "Seller Update Request"
// @Success 200 {object} SellersGetJSON "Updated Seller data"
// @Failure 400 {object} resterr.RestErr "Seller invalid fields"
// @Failure 404 {object} resterr.RestErr "Seller not found"
// @Failure 409 {object} resterr.RestErr "Seller with this CID already exists"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/sellers/{id} [patch]
func (h *SellerDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
			return
		}

		var body SellersUpdateJSON

		err = request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError("request json invalid"))
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
			h.handleError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": SellersGetJSON{
				ID:          seller.ID,
				CID:         seller.CID,
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
// @Failure 400 {object} resterr.RestErr "Bad Request"
// @Failure 404 {object} resterr.RestErr "Seller not found"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/sellers/{id} [delete]
func (h *SellerDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
			return
		}

		err = h.sv.Delete(id)
		if err != nil {
			h.handleError(w, err)
			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}

func (h *SellerDefault) handleError(w http.ResponseWriter, err error) {
	if errors.As(err, &internal.DomainError{}) {
		var domainError internal.DomainError
		_ = errors.As(err, &domainError)

		var causes []resterr.Causes
		for _, cause := range domainError.Causes {
			causes = append(causes, resterr.Causes{
				Field:   cause.Field,
				Message: cause.Message,
			})
		}

		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestValidationError(domainError.Message, causes))

		return
	}

	if errors.Is(err, internal.ErrSellerInvalidFields) {
		restErr := resterr.NewBadRequestError(err.Error())
		response.JSON(w, restErr.Code, restErr)

		return
	}

	if errors.Is(err, internal.ErrSellerConflict) || errors.Is(err, internal.ErrSellerCIDAlreadyExists) {
		response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))

		return
	}

	if errors.Is(err, internal.ErrSellerNotFound) || errors.Is(err, internal.ErrLocalityNotFound) {
		response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))

		return
	}

	response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError("internal server error"))
}
