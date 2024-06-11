// Package coupon contain all domain logic & needed interfaces
package coupon

import (
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
)

var (
	// ErrCouponAlreadyUsed used when the coupon is already used
	ErrCouponAlreadyUsed = internalErrors.NewConflict("coupon already used")
	// ErrCouponEmptyName used when coupon has empty name
	ErrCouponEmptyName = internalErrors.NewWrongInput("coupon empty name")
	// ErrCouponInvalidAmount used when coupon has invalid amount
	ErrCouponInvalidAmount = internalErrors.NewWrongInput("coupon invalid amount")
)

// Coupon defines the asset of a coupon in our service
type Coupon struct {
	// ID Unique Identifier of the Coupon
	ID uuid.UUID `json:"id,omitempty"`
	// Name will be the name of the Coupon
	Name string `json:"name,omitempty"`
	// Amount that will be used to deduct from shopping cart
	Amount float32 `json:"amount,omitempty"`
	// Used will represent if the coupon had been used
	Used bool `json:"used,omitempty"`
	// Timestamp when it was created
	CreatedAt time.Time `json:"created_at,omitempty"`
	// Timestamp of the last update
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// New return a new Coupon instance
func New(req CreateRequest) *Coupon {
	return &Coupon{
		ID:     uuid.MustParse(uuid.NewString()),
		Name:   req.Name,
		Amount: float32(math.Round(float64(req.Amount))),
		Used:   false,
	}
}

// IsUsed checks if coupon can be used
func (c *Coupon) IsUsed() bool {
	return c.Used
}

// CreateRequest defines needed field to create a coupon
type CreateRequest struct {
	Name   string  `json:"name,omitempty"`
	Amount float32 `json:"amount,omitempty"`
}

// Validate validates the create request
func (r CreateRequest) Validate() error {
	if r.Name == "" {
		return ErrCouponEmptyName
	}
	if r.Amount <= 0 {
		return ErrCouponInvalidAmount
	}
	return nil
}

// Service defines the available functions for the Coupon Service
type Service interface {
	// CreateCoupon returns a new coupon
	CreateCoupon(CreateRequest) (*Coupon, error)
	// ListCoupons returns a list of coupons
	ListCoupons() ([]Coupon, error)
}

// Repository defines the available functions for the Coupon repository
type Repository interface {
	// CreateCoupon returns a new coupon
	CreateCoupon(*Coupon) (*Coupon, error)
	// ListCoupons returns a list of coupons
	ListCoupons() ([]Coupon, error)
	// GetCouponForUpdate returns an specific  and it will lock the row in order to update it
	GetCouponForUpdate(*gorm.DB, uuid.UUID) (*Coupon, error)
	// UpdateCoupon updates coupon entity
	UpdateCoupon(*gorm.DB, *Coupon) (*Coupon, error)
}

// Server defines what are the different allowed http
// endpoints that can be consumed
type Server interface {
	// CreateCoupon returns a new coupon
	CreateCoupon(w http.ResponseWriter, r *http.Request)
	// ListCoupons returns a list of coupons
	LisCoupons(w http.ResponseWriter, r *http.Request)
}
