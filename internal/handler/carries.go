package handler

import (
	"log"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type CarriesHandlerDefault struct {
	sv internal.CarriesService
}

func NewCarriesHandlerDefault(sv internal.CarriesService) *CarriesHandlerDefault {
	return &CarriesHandlerDefault{sv}
}

func (h *CarriesHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	all, err := h.sv.FindAll()
	if err != nil {
		log.Println(err)
		response.JSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to fetch carries",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": all,
	})
}
