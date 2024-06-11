package service

import (
	"github.com/google/uuid"

	couponDomain "github.com/nachoconques0/schwarz-challenge/internal/coupon"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
)

type shoppingCartService struct {
	shoppingCartRepo shoppingcart.Repository
	couponRepo       couponDomain.Repository
}

// NewShoppingCartRepository builds a new repository that
// satisfies the shopping cart interface
func NewShoppingCartService(scr shoppingcart.Repository, cr couponDomain.Repository) (shoppingcart.Service, error) {
	return &shoppingCartService{
		shoppingCartRepo: scr,
		couponRepo:       cr,
	}, nil
}

// CreateShoppingCart will create a new shopping cart
func (sc *shoppingCartService) CreateShoppingCart(req shoppingcart.CreateRequest) (*shoppingcart.ShoppingCart, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	payload := shoppingcart.New(req)
	res, err := sc.shoppingCartRepo.CreateShoppingCart(payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ListShoppingCarts returns a shopping cart list
func (sc *shoppingCartService) ListShoppingCarts() ([]shoppingcart.ShoppingCart, error) {
	res, err := sc.shoppingCartRepo.ListShoppingCarts()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ApplyCoupon applies a coupon code
func (sc *shoppingCartService) ApplyCoupon(scID uuid.UUID, couponID uuid.UUID) error {
	tx := sc.shoppingCartRepo.BeginTransaction()

	coupon, err := sc.couponRepo.GetCouponForUpdate(tx, couponID)
	if err != nil {
		return err
	}

	if coupon.IsUsed() {
		return couponDomain.ErrCouponAlreadyUsed
	}

	toUpdateShoppingCart, err := sc.shoppingCartRepo.GetShoppingCartForUpdate(tx, scID)
	if err != nil {
		return err
	}

	err = toUpdateShoppingCart.ApplyCoupon(couponID, coupon.Amount)
	if err != nil {
		return err
	}

	_, err = sc.shoppingCartRepo.UpdateShoppingCart(tx, toUpdateShoppingCart)
	if err != nil {
		return err
	}

	coupon.Used = true
	_, err = sc.couponRepo.UpdateCoupon(tx, coupon)
	if err != nil {
		return err
	}
	err = sc.shoppingCartRepo.CommitTransaction(tx)
	if err != nil {
		return err
	}
	return nil
}
