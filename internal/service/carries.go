package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

type CarriesService struct {
	rp internal.CarriesRepository
}

func NewCarriesService(rp internal.CarriesRepository) *CarriesService {
	return &CarriesService{rp}
}

func (sv *CarriesService) FindAll() ([]internal.Carries, error) {
	return sv.rp.FindAll()
}
