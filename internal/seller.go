package internal

import "errors"

// Seller is a struct that contains the seller's information
type Seller struct {
	// ID is the unique identifier of the seller
	ID int
	// CID is the unique identifier of the company
	CID int
	// CompanyName is the name of the company
	CompanyName string
	// Address is the address of the company
	Address string
	// Telephone is the telephone number of the company
	Telephone string
}

var (
	// ErrSellerNotFound is returned when the seller is not found
	ErrSellerNotFound = errors.New("seller not found")
	// ErrSellerConflict is returned when the seller already exists
	ErrSellerConflict = errors.New("seller already exists")
)

// SellerRepository is an interface that contains the methods that the seller repository should support
type SellerRepository interface {
	// FindAll returns all the sellers
	FindAll() ([]Seller, error)
	// FindByID returns the seller with the given ID
	FindByID(id int) (Seller, error)
	// Save saves the given seller
	Save(seller *Seller) error
	// Update updates the given seller
	Update(seller *Seller) error
	// Delete deletes the seller with the given ID
	Delete(id int) error
}

// SellerService is an interface that contains the methods that the seller service should support
type SellerService interface {
	// FindAll returns all the sellers
	FindAll() ([]Seller, error)
	// FindByID returns the seller with the given ID
	FindByID(id int) (Seller, error)
	// Save saves the given seller
	Save(seller *Seller) error
	// Update updates the given seller
	Update(seller *Seller) error
	// Delete deletes the seller with the given ID
	Delete(id int) error
}
