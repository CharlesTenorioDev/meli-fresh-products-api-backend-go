package service

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type LocalityDefault struct {
	rp internal.LocalityRepository
}

func NewLocalityDefault(rp internal.LocalityRepository) *LocalityDefault {
	return &LocalityDefault{
		rp,
	}
}

func (l *LocalityDefault) Save(locality *internal.Locality) (err error) {
	causes := locality.Validate()

	if len(causes) > 0 {
		return internal.DomainError{
			Message: "Locality inputs are Invalid",
			Causes:  causes,
		}
	}

	return l.rp.Save(locality)
}

func (l *LocalityDefault) ReportSellers() (localities []internal.Locality, err error) {
	return l.rp.ReportSellers()
}
func (l *LocalityDefault) ReportSellersByID(id int) (localities []internal.Locality, err error) {
	return l.rp.ReportSellersByID(id)
}

func (l *LocalityDefault) FindByID(id int) (locality internal.Locality, err error) {
	return l.rp.FindByID(id)
}

func (sv *LocalityDefault) ReportCarries(localityId int) (int, error) {
	return sv.rp.ReportCarries(localityId)
}

func (sv *LocalityDefault) GetAmountOfCarriesForEveryLocality() ([]internal.CarriesCountPerLocality, error) {
	return sv.rp.GetAmountOfCarriesForEveryLocality()
}
