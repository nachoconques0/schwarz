package repo

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// ErrCouponNotFound used when coupon is not found
	ErrCouponNotFound = internalErrors.NewNotFound("coupon not found")
	// ErrCouponMissingID used when coupon id is missing
	ErrCouponMissingID = internalErrors.NewWrongInput("coupon id is missing")
)

// couponTable is the table name for the coupon model
const couponTable = "schwarz.coupon"

type couponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) (coupon.Repository, error) {
	if db == nil {
		return nil, ErrMissingDB
	}
	return &couponRepository{
		db: db,
	}, nil
}

// CreateCoupon returns a new coupon
func (cs couponRepository) CreateCoupon(coupon *coupon.Coupon) (*coupon.Coupon, error) {
	if err := cs.db.Table(couponTable).Create(&coupon).Error; err != nil {
		return nil, err
	}
	return coupon, nil
}

// ListCoupons returns a list of coupons
func (cs couponRepository) ListCoupons() ([]coupon.Coupon, error) {
	var result []coupon.Coupon
	if err := cs.db.Table(couponTable).Find(&result).Error; err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrCouponNotFound
	}
	return result, nil
}

// GetCouponForUpdate returns an specific  and it will lock the row in order to update it
func (cs couponRepository) GetCouponForUpdate(tx *gorm.DB, couponID uuid.UUID) (*coupon.Coupon, error) {
	if couponID == uuid.Nil {
		return nil, ErrCouponMissingID
	}

	var result *coupon.Coupon
	if err := tx.Table(couponTable).
		Where(coupon.Coupon{ID: couponID}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCouponNotFound
		}
		return nil, err
	}
	return result, nil
}

// UpdateCoupon updates coupon entity
func (cs couponRepository) UpdateCoupon(tx *gorm.DB, coupon *coupon.Coupon) (*coupon.Coupon, error) {
	if err := tx.Table(couponTable).Save(&coupon).Error; err != nil {
		return nil, err
	}
	return coupon, nil
}
