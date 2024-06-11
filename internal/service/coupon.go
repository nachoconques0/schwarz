package service

import (
	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/errors"
)

var (
	// ErrMissingDB used when DB is nil
	ErrMissingDB = errors.NewNotFound("DB connection is missing")
)

type couponService struct {
	repo coupon.Repository
}

// NewCouponService builds a new repository that
// satisfies the coupon interface
func NewCouponService(repo coupon.Repository) (coupon.Service, error) {
	return &couponService{
		repo: repo,
	}, nil
}

// CreateCoupon creates a new coupon
func (cs *couponService) CreateCoupon(req coupon.CreateRequest) (*coupon.Coupon, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}
	payload := coupon.New(req)
	res, err := cs.repo.CreateCoupon(payload)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListCoupons returns a list of coupons
func (cs *couponService) ListCoupons() ([]coupon.Coupon, error) {
	res, err := cs.repo.ListCoupons()
	if err != nil {
		return nil, err
	}
	return res, nil
}
