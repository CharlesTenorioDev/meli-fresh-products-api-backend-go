package internal

type Buyer struct {
	ID           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type BuyerRepository interface {
	GetAll() (db map[int]Buyer)
	AddProduct(id int, buyer Buyer)
}

type BuyerService interface {
	GetAll() (map[int]Buyer)
	FindByID(id int) (b Buyer, err error)
	Save(id int, buyer Buyer) (err error)
}

func (b *Buyer) Parse() (ok bool) {
	ok = true
	if b.CardNumberId == "" || b.LastName == "" || b.FirstName == "" {
		ok = false
	}
	return
}
