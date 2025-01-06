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
	Create(carry Carries) (lastId int64, e error)
}

type CarriesRepository interface {
	FindAll() ([]Carries, error)
	Create(carry Carries) (lastId int64, e error)
}

func (c *Carries) Ok() bool {
	if c.Cid < 0 || c.CompanyName == "" || c.Address == "" || c.PhoneNumber == "" || c.LocalityId < 0 {
		return false
	}
	return true
}
