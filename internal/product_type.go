package internal

type ProductType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProductTypeRepository interface {
	FindByID(id int) (ProductType, error)
}
