package internal

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/utils/validator"
)

type Locality struct {
	ID           int
	LocalityName string
	ProvinceName string
	CountryName  string
	Sellers      int
}

type CarriesCountPerLocality struct {
	CarriesCount int    `json:"carries_count"`
	LocalityID   int    `json:"locality_id"`
	LocalityName string `json:"locality_name"`
}

var (
	// ErrLocalityNotFound is returned when the seller is not found
	ErrLocalityNotFound = errors.New("locality not found")
	ErrLocalityConflict = errors.New("locality conflict")
)

func (l *Locality) Validate() (causes []Causes) {
	if validator.IntIsNegative(l.ID) {
		causes = append(causes, Causes{
			Field:   "locality_id",
			Message: "Locality ID cannot be negative",
		})
	}

	if validator.IntIsZero(l.ID) {
		causes = append(causes, Causes{
			Field:   "locality_id",
			Message: "Locality ID is required",
		})
	}

	if !validator.String(l.LocalityName, 1, 255) {
		causes = append(causes, Causes{
			Field:   "locality_name",
			Message: "Locality name is required",
		})
	}

	if !validator.String(l.CountryName, 1, 255) {
		causes = append(causes, Causes{
			Field:   "country_name",
			Message: "Country name cannot be empty",
		})
	}

	if !validator.String(l.ProvinceName, 1, 255) {
		causes = append(causes, Causes{
			Field:   "province_name",
			Message: "Province name cannot be empty",
		})
	}

	return causes
}

type LocalityRepository interface {
	Save(locality *Locality) (err error)
	ReportSellers() (localities []Locality, err error)
	ReportSellersByID(id int) (localities []Locality, err error)
	FindByID(id int) (locality Locality, err error)
	ReportCarries(localityID int) (amountOfCarries int, e error)
	GetAmountOfCarriesForEveryLocality() (c []CarriesCountPerLocality, e error)
}

type LocalityService interface {
	Save(locality *Locality) (err error)
	ReportSellers() (localities []Locality, err error)
	ReportSellersByID(id int) (localities []Locality, err error)
	FindByID(id int) (locality Locality, err error)
	ReportCarries(localityID int) (int, error)
	GetAmountOfCarriesForEveryLocality() ([]CarriesCountPerLocality, error)
}
