package handler

import (
	"errors"
	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
	"net/http"
	"strconv"
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
	Sellers      int    `json:"report_sellers"`
}

// ReportSellers returns locality with sellers count
func (h *LocalityDefault) ReportSellers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")

		id, err := strconv.Atoi(idStr)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, nil)
			return
		}

		locality, err := h.sv.ReportSellers(id)
		if err != nil {
			if errors.Is(err, internal.ErrLocalityNotFound) {
				response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
				return
			}
			response.JSON(w, http.StatusInternalServerError, nil)
			return
		}

		localityJson := LocalityGetJson{
			ID:           locality.ID,
			LocalityName: locality.LocalityName,
			ProvinceName: locality.ProvinceName,
			CountryName:  locality.CountryName,
			Sellers:      locality.Sellers,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": localityJson,
		})
	}
}
