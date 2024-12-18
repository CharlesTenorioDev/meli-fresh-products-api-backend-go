package internal

type Employee struct {
	Id           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

type EmployeePatch struct {
	CardNumberId *string `json:"card_number_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
	WarehouseId  *int    `json:"warehouse_id,omitempty"`
}

func (emp *Employee) RequirementsFields() (ok bool) {
	ok = true

	if emp.CardNumberId == "" || emp.FirstName == "" || emp.LastName == "" || emp.WarehouseId == 0 {
		ok = false
	}
	return
}
