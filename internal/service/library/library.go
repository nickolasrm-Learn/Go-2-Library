package library

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nickolasrm-Learn/Go-2-Library/internal/model/product"
	"github.com/nickolasrm-Learn/Go-2-Library/internal/model/purchase"
	"github.com/nickolasrm-Learn/Go-2-Library/internal/model/user"
)

// Library represents a collection of users, products and purchases
type Library struct {
	Users     map[string]*user.User
	Products  map[string]*product.Product
	Purchases map[string]*purchase.Purchase
}

// LibraryPath is the path to the library file
var LibraryPath = "library.json"

// Load a library object from a JSON file
//
// Returns nil if the file doesn't exist
// Returns an error if the file exists but it's not a valid JSON
func Load() (*Library, error) {
	content, err := os.ReadFile(LibraryPath)
	if err != nil {
		return nil, nil
	}
	var container Library
	err = json.Unmarshal(content, &container)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// Save a library object to a JSON file
//
// Returns an error if the file can't be written
func (l *Library) Save() error {
	jsonString, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(LibraryPath, jsonString, 0644)
	if err != nil {
		return err
	}
	return nil
}

// New creates a new library object
//
// If a library file exists, it loads it
// Otherwise, it creates a new empty library
func New() (*Library, error) {
	library, err := Load()
	if err != nil {
		return nil, err
	}
	if library == nil {
		return &Library{
			Users:     make(map[string]*user.User),
			Products:  make(map[string]*product.Product),
			Purchases: make(map[string]*purchase.Purchase),
		}, nil
	}
	return library, nil
}

// AddUser adds a new user to the library
func (l *Library) AddUser(name string) *user.User {
	usr := user.New(name)
	l.Users[usr.ID] = usr
	return usr
}

// AddProduct adds a new product to the library
func (l *Library) AddProduct(title string, price float64, category product.Category) (*product.Product, error) {
	prod, err := product.New(title, price, category)
	if err != nil {
		return nil, err
	}
	l.Products[prod.ID] = prod
	return prod, nil
}

// ListProducts returns a list of all products in the library
func (l *Library) ListProducts() []*product.Product {
	lst := make([]*product.Product, len(l.Products))
	i := 0
	for _, v := range l.Products {
		lst[i] = v
		i++
	}
	return lst
}

// ListPurchases returns a list of all purchases made by a user
func (l *Library) ListPurchases(usr *user.User) []*purchase.Purchase {
	lst := make([]*purchase.Purchase, 0)
	for _, v := range l.Purchases {
		if v.User == usr.ID {
			lst = append(lst, v)
		}
	}
	return lst
}

// GetUser returns a user by its ID
func (l *Library) GetUser(userID string) (*user.User, error) {
	v, ok := l.Users[userID]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("user not found")
}

// GetProduct returns a product by its ID
func (l *Library) GetProduct(productID string) (*product.Product, error) {
	v, ok := l.Products[productID]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("product not found")
}

// Deposit adds value to user budget
func (l *Library) Deposit(usr *user.User, value float64) {
	usr.Deposit(value)
}

// Buy creates a new purchase from user and product
//
// Returns an error if user doesn't have enough budget
func (l *Library) Buy(usr *user.User, prod *product.Product) (*purchase.Purchase, error) {
	err := usr.Spend(prod.Price)
	if err != nil {
		return nil, err
	}
	prc := purchase.New(usr, prod)
	l.Purchases[prc.ID] = prc
	return prc, nil
}
