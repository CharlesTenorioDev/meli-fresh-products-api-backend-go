package internal

type Buyer struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type BuyerPatch struct {
	CardNumberID *string `json:"card_number_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
}

type PurchaseOrdersByBuyer struct {
	BuyerID             int    `json:"id"`
	CardNumberID        string `json:"card_number_id"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	PurchaseOrdersCount int    `json:"purchase_orders_count"`
}

type BuyerRepository interface {
	GetAll() (db map[int]Buyer, err error)
	Add(buyer *Buyer) (id int64, err error)
	Update(id int, buyer BuyerPatch) (err error)
	Delete(id int) (rowsAffected int64, err error)
	ReportPurchaseOrders() (purchaseOrders []PurchaseOrdersByBuyer, err error)
	ReportPurchaseOrdersByID(id int) (purchaseOrders []PurchaseOrdersByBuyer, err error)
}

type BuyerService interface {
	GetAll() map[int]Buyer
	FindByID(id int) (b Buyer, err error)
	Save(buyer *Buyer) (err error)
	Update(id int, buyerPatch BuyerPatch) (err error)
	Delete(id int) (err error)
	ReportPurchaseOrders() (purchaseOrders []PurchaseOrdersByBuyer, err error)
	ReportPurchaseOrdersByID(id int) (purchaseOrders []PurchaseOrdersByBuyer, err error)
}

func (b *Buyer) Parse() (ok bool) {
	ok = true
	if b.CardNumberID == "" || b.LastName == "" || b.FirstName == "" {
		ok = false
	}

	return
}

func (b BuyerPatch) Patch(buyerToUpdate *Buyer) {
	if b.CardNumberID != nil {
		buyerToUpdate.CardNumberID = *b.CardNumberID
	}

	if b.FirstName != nil {
		buyerToUpdate.FirstName = *b.FirstName
	}

	if b.LastName != nil {
		buyerToUpdate.LastName = *b.LastName
	}
}
