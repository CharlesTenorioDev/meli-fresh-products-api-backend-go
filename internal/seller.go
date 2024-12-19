package internal

import (
	"errors"
)

// Seller is a struct that contains the seller's information
type Seller struct {
	// ID is the unique identifier of the seller
	ID int `json:"id"`
	// CID is the unique identifier of the company
	CID int `json:"cid"`
	// CompanyName is the name of the company
	CompanyName string `json:"company_name"`
	// Address is the address of the company
	Address string `json:"address"`
	// Telephone is the telephone number of the company
	Telephone string `json:"telephone"`
}

type SellerPatch struct {
	CID         *int
	CompanyName *string
	Address     *string
	Telephone   *string
}

func (seller *Seller) Validate() error {
	var err error
	if seller.CID == 0 {
		return errors.Join(err, errors.New("seller.CID is required"))
	}
	if seller.CompanyName == "" {
		return errors.Join(err, errors.New("seller.CompanyName is required"))
	}
	if seller.Address == "" {
		return errors.Join(err, errors.New("seller.Address is required"))
	}
	if seller.Telephone == "" {
		return errors.Join(err, errors.New("seller.Telephone is required"))
	}
	return err
}

var (
	ErrSellerCIDAlreadyExists = errors.New("seller with this CID already exists")
	ErrSellerInvalidFields    = errors.New("seller invalid fields")
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
	// FindByCID returns the seller with the given CID
	FindByCID(cid int) (*Seller, error)
	// Save saves the given seller
	Save(seller *Seller) error
	// Update updates the given seller
	Update(id int, seller *Seller) error
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
	Update(id int, updateSeller SellerPatch) (Seller, error)
	// Delete deletes the seller with the given ID
	Delete(id int) error
}
