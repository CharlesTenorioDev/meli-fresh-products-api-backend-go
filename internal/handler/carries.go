package handler

import (
	"encoding/json"
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

func (h *CarriesHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var carry internal.Carries
	err := json.NewDecoder(r.Body).Decode(&carry)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"error": "failed to parse body",
		})
		return
	}

	if ok := carry.Ok(); !ok {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "missing fields",
		})
		return
	}

	lastId, err := h.sv.Create(carry)
	if err != nil {
		response.JSON(w, http.StatusConflict, map[string]any{
			"error": "conflict",
		})
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": struct {
			Id int64 `json:"id"`
		}{
			Id: lastId,
		},
	})
}
