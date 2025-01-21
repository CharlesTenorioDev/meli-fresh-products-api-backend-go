package internal

type Carries struct {
	ID          int    `json:"id"`
	Cid         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	LocalityID  int    `json:"locality_id"`
}

type CarriesService interface {
	FindAll() (carries []Carries, e error)
	Create(carry Carries) (lastID int64, e error)
}

type CarriesRepository interface {
	FindAll() ([]Carries, error)
	Create(carry Carries) (lastID int64, e error)
}

func (c *Carries) Ok() bool {
	if c.Cid == "" || c.CompanyName == "" || c.Address == "" || c.PhoneNumber == "" || c.LocalityID < 0 {
		return false
	}

	return true
}
