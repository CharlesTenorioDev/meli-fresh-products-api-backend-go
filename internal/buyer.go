package internal

type Buyer struct {
	ID           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type BuyerPatch struct {
	CardNumberId *string `json:"card_number_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
}

type BuyerRepository interface {
	GetAll() (db map[int]Buyer)
	AddProduct(id int, buyer Buyer)
	UpdateBuyer(id int, buyer BuyerPatch)
}

type BuyerService interface {
	GetAll() (map[int]Buyer)
	FindByID(id int) (b Buyer, err error)
	Save(id int, buyer Buyer) (err error)
	Update(id int, buyerPatch BuyerPatch) (err error)
}

func (b *Buyer) Parse() (ok bool) {
	ok = true
	if b.CardNumberId == "" || b.LastName == "" || b.FirstName == "" {
		ok = false
	}
	return
}

func (b BuyerPatch) Patch(buyerToUpdate *Buyer) {
	if b.CardNumberId != nil {
		buyerToUpdate.CardNumberId = *b.CardNumberId
	}

	if b.FirstName != nil {
		buyerToUpdate.FirstName = *b.FirstName
	}

	if b.LastName != nil {
		buyerToUpdate.LastName = *b.LastName
	}
}
