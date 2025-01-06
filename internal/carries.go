package internal

type Carries struct {
	Id          int    `json:"id"`
	Cid         int    `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	LocalityId  int    `json:"locality_id"`
}

type CarriesService interface {
	FindAll() (carries []Carries, e error)
}

type CarriesRepository interface {
	FindAll() ([]Carries, error)
}
