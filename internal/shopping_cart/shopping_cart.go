// Package shoppingcart contain all domain logic & needed interfaces
package shoppingcart

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
)

var (
	// ErrShoppingCartWithCouponAlreadyApplied used when shopping cart has a coupon applied
	ErrShoppingCartWithCouponAlreadyApplied = internalErrors.NewConflict("shopping cart coupon already used")
	// ErrShoppingCartEmptyItems used when items are empty
	ErrShoppingCartEmptyItems = internalErrors.NewWrongInput("shopping cart empty items")
	// ErrItemEmptyName used when item has empty name
	ErrItemEmptyName = internalErrors.NewWrongInput("item empty name")
	// ErrItemEmptyDescription used when item has empty description
	ErrItemEmptyDescription = internalErrors.NewWrongInput("item empty description")
	// ErrItemEmptyInvalidAmount used when item has invalid amount
	ErrItemEmptyInvalidAmount = internalErrors.NewWrongInput("item invalid amount")
	// ErrShoppinCartCouponAlreadyApplied used when a coupon was already applied
	ErrShoppinCartCouponAlreadyApplied = internalErrors.NewConflict("shopping cart already with coupon applied")
	// ErrShoppointCartCouponAmountExceeded used when a coupon amount was
	ErrShoppointCartCouponAmountExceeded = internalErrors.NewWrongInput("coupon amount exceeds shopping cart total")
)

// ShoppingCart defines the asset of a Shopping Cart in our service
type ShoppingCart struct {
	// ID Unique Identifier of the Shopping Cart
	ID uuid.UUID `json:"id,omitempty"`
	// Items will be an array of Items associated to the shopping cart
	Items Items `json:"items,omitempty"`
	// Amount is the total before any discounts applied
	Amount float32 `json:"amount,omitempty"`
	// Total is the aggregate amount of the amount and discounts (if applied)
	Total float32 `json:"total,omitempty"`
	// CouponID will be the ID of the applied coupon
	CouponID uuid.UUID `json:"coupon_id,omitempty"`
	// Timestamp when it was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Timestamp of the last update
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// CreateRequest defines needed field to create a shopping cart
type CreateRequest struct {
	Items Items `json:"items,omitempty"`
}

// Item defines the asset of a Item in our service
type Item struct {
	// ID Unique Identifier of an Item
	ID uuid.UUID `json:"id,omitempty"`
	// Name will be the name of the item
	Name string `json:"name,omitempty"`
	// Description will be the description of the item
	Description string `json:"description,omitempty"`
	// Price
	Price float32 `json:"price,omitempty"`
}

// Validate validates the create request
func (r CreateRequest) Validate() error {
	if len(r.Items) == 0 {
		return ErrShoppingCartEmptyItems
	}
	for _, i := range r.Items {
		err := i.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

// Validate validates the item fields
func (i Item) Validate() error {
	if i.Name == "" {
		return ErrItemEmptyName
	}
	if i.Description == "" {
		return ErrItemEmptyDescription
	}
	if i.Price <= 0 {
		return ErrItemEmptyInvalidAmount
	}
	return nil
}

// New returns a new Shopping Cart instance
func New(req CreateRequest) *ShoppingCart {
	var parsedItems []Item
	var totalAmount float32
	for _, i := range req.Items {
		totalAmount += i.Price
		parsedItems = append(parsedItems, Item{
			ID:          uuid.MustParse(uuid.NewString()),
			Name:        i.Name,
			Description: i.Description,
			Price:       i.Price,
		})
	}
	return &ShoppingCart{
		ID:     uuid.MustParse(uuid.NewString()),
		Items:  parsedItems,
		Amount: float32(int(totalAmount*100)) / 100,
		Total:  float32(int(totalAmount*100)) / 100,
	}
}

// ApplyCoupon will deduct the coupon amount from the total of the shopping cart
func (sc *ShoppingCart) ApplyCoupon(couponID uuid.UUID, couponAmount float32) error {
	if sc.CouponID != uuid.Nil {
		return ErrShoppinCartCouponAlreadyApplied
	}
	if couponAmount >= sc.Total {
		return ErrShoppointCartCouponAmountExceeded
	}

	sc.CouponID = couponID
	res := sc.Total - couponAmount
	sc.Total = float32(int(res*100)) / 100

	return nil
}

// Items contains list items in json format
type Items []Item

// Value for DB
func (i Items) Value() (driver.Value, error) {
	res, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Scan will unmarshall Items data
func (i *Items) Scan(src interface{}) error {
	switch t := src.(type) {
	case string:
		err := json.Unmarshal([]byte(t), &i)
		if err != nil {
			return err
		}
		return nil
	case []byte:
		err := json.Unmarshal(t, &i)
		if err != nil {
			return err
		}
		return nil
	case nil:
		*i = nil
		return nil
	}
	return errors.New("err unmarshal entity")
}

// Service defines the available functions for the Shopping Cart Service
type Service interface {
	// CreateShoppingCart will create a new shopping cart
	CreateShoppingCart(CreateRequest) (*ShoppingCart, error)
	// ListShoppingCarts returns a shopping cart list
	ListShoppingCarts() ([]ShoppingCart, error)
	// ApplyCoupon applies a coupon code
	ApplyCoupon(uuid.UUID, uuid.UUID) error
}

// Repository defines the available functions for the Shopping Cart repository
type Repository interface {
	// CreateShoppingCart will create a new shopping cart
	CreateShoppingCart(*ShoppingCart) (*ShoppingCart, error)
	// GetShoppingCartForUpdate returns a shopping cart and it will lock the row in order to update it
	GetShoppingCartForUpdate(*gorm.DB, uuid.UUID) (*ShoppingCart, error)
	// ListShoppingCarts returns a list of shopping carts
	ListShoppingCarts() ([]ShoppingCart, error)
	// UpdateShoppingCart updates shopping cart entity
	UpdateShoppingCart(*gorm.DB, *ShoppingCart) (*ShoppingCart, error)
	BeginTransaction() *gorm.DB
	CommitTransaction(tx *gorm.DB) error
	RollbackTransaction(tx *gorm.DB) error
}

// Server defines what are the different allowed http
// endpoints that can be consumed
type Server interface {
	// CreateShoppingCart receives a request in order to create a shopping cart
	CreateShoppingCart(w http.ResponseWriter, r *http.Request)
	// ListShoppingCarts returns a list of shopping carts
	ListShoppingCarts(w http.ResponseWriter, r *http.Request)
	// ApplyCoupon receives a request in order to apply a coupon to a shopping cart
	ApplyCoupon(w http.ResponseWriter, r *http.Request)
}
