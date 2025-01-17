package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/logger"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

// NewLocalityDefault creates a new instance of the seller handler
func NewLocalityDefault(sv internal.LocalityService) *LocalityDefault {
	return &LocalityDefault{
		sv: sv,
	}
}

// LocalityDefault is the default implementation of the seller handler
type LocalityDefault struct {
	// sv is the service used by the handler
	sv internal.LocalityService
}

type LocalityGetJson struct {
	ID           int    `json:"id"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
	SellersCount int    `json:"sellers_count"`
}

type LocalityPostJson struct {
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}

func (h *LocalityDefault) ReportCarries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")

		if idStr == "" {
			carries, err := h.sv.GetAmountOfCarriesForEveryLocality()
			if err != nil {
				response.JSON(
					w,
					http.StatusInternalServerError,
					resterr.NewInternalServerError("failed to fetch carries"),
				)
				return
			}

			response.JSON(w, http.StatusOK, map[string]any{
				"data": carries,
			})
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(
				w,
				http.StatusBadRequest,
				resterr.NewBadRequestError("id should be a number"),
			)
			return
		}

		amountOfCarries, err := h.sv.ReportCarries(id)
		if err != nil {
			response.JSON(
				w,
				http.StatusNotFound,
				resterr.NewNotFoundError("not carries on locality_id "+idStr),
			)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": struct {
				AmountOfCarries int `json:"amount_of_carries"`
			}{
				AmountOfCarries: amountOfCarries,
			},
		})
	}
}

// ReportSellers returns locality with sellers count
func (h *LocalityDefault) ReportSellers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var localities []internal.Locality
		var err error

		idStr := r.URL.Query().Get("id")

		switch idStr {
		case "":
			localities, err = h.sv.ReportSellers()
		default:
			id, parseErr := strconv.Atoi(idStr)

			if parseErr != nil {
				response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError("id should be a number"))
				return
			}

			localities, err = h.sv.ReportSellersByID(id)
		}

		if err != nil {
			logger.Error(err.Error(), err,
				zap.String("id", idStr),
			)
			if errors.Is(err, internal.ErrLocalityNotFound) {
				response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
				return
			}
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		var localitiesJson []LocalityGetJson
		for _, locality := range localities {
			localitiesJson = append(localitiesJson, LocalityGetJson{
				ID:           locality.ID,
				LocalityName: locality.LocalityName,
				ProvinceName: locality.ProvinceName,
				CountryName:  locality.CountryName,
				SellersCount: locality.Sellers,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": localitiesJson,
		})
	}

}

// Save method save the locality
func (h *LocalityDefault) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var localityJson LocalityPostJson
		err := json.NewDecoder(r.Body).Decode(&localityJson)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		locality := &internal.Locality{
			ID:           localityJson.LocalityID,
			LocalityName: localityJson.LocalityName,
			ProvinceName: localityJson.ProvinceName,
			CountryName:  localityJson.CountryName,
		}

		err = h.sv.Save(locality)
		if err != nil {
			if errors.Is(err, internal.ErrLocalityConflict) {
				response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
				return
			}

			if errors.As(err, &internal.DomainError{}) {
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
				return
			}

			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": localityJson,
		})
	}
}
