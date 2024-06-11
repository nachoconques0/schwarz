package coupon_test

import (
	"testing"

	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/stretchr/testify/assert"
)

var (
	testName   = "name"
	testAmount = 100
)

func TestCouponNew(t *testing.T) {
	c := coupon.New(coupon.CreateRequest{
		Name:   testName,
		Amount: float32(testAmount),
	})
	assert.Equal(t, c.Amount, float32(testAmount))
	assert.Equal(t, c.Name, testName)
}

func TestCouponCreateValidate(t *testing.T) {
	req := coupon.CreateRequest{
		Name:   testName,
		Amount: float32(testAmount),
	}

	t.Run("invalid name", func(t *testing.T) {
		req.Name = ""
		err := req.Validate()
		assert.Equal(t, coupon.ErrCouponEmptyName, err)
	})
	t.Run("invalid amount", func(t *testing.T) {
		req.Amount = 0
		req.Name = "ol"
		err := req.Validate()
		assert.Equal(t, coupon.ErrCouponInvalidAmount, err)
	})
}
