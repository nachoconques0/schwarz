package shoppingcart_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
	"github.com/stretchr/testify/assert"
)

var (
	testName        = "name"
	testDescription = "description"
	testAmount      = 100
)

func TestShoppingCartNew(t *testing.T) {
	sc := shoppingcart.New(shoppingcart.CreateRequest{
		Items: shoppingcart.Items{
			shoppingcart.Item{
				Name:  testName,
				Price: float32(testAmount),
			},
		},
	})
	assert.Len(t, sc.Items, 1)
	assert.Equal(t, sc.Amount, float32(testAmount))
	assert.Equal(t, sc.Items[0].Name, testName)
}

func TestShoppingCartCreateValidate(t *testing.T) {
	t.Run("invalid items", func(t *testing.T) {
		req := shoppingcart.CreateRequest{
			Items: shoppingcart.Items{},
		}
		err := req.Validate()
		assert.Equal(t, shoppingcart.ErrShoppingCartEmptyItems, err)
	})

	t.Run("invalid item name", func(t *testing.T) {
		req := shoppingcart.CreateRequest{
			Items: shoppingcart.Items{
				shoppingcart.Item{
					Price:       float32(testAmount),
					Description: testDescription,
				},
			},
		}
		err := req.Validate()
		assert.Equal(t, shoppingcart.ErrItemEmptyName, err)
	})

	t.Run("invalid item description", func(t *testing.T) {
		req := shoppingcart.CreateRequest{
			Items: shoppingcart.Items{
				shoppingcart.Item{
					Price: float32(testAmount),
					Name:  testName,
				},
			},
		}
		err := req.Validate()
		assert.Equal(t, shoppingcart.ErrItemEmptyDescription, err)
	})

	t.Run("invalid item price", func(t *testing.T) {
		req := shoppingcart.CreateRequest{
			Items: shoppingcart.Items{
				shoppingcart.Item{
					Description: testDescription,
					Name:        testName,
				},
			},
		}
		err := req.Validate()
		assert.Equal(t, shoppingcart.ErrItemEmptyInvalidAmount, err)
	})
}

func TestShoppingCartApplyCoupon(t *testing.T) {
	req := shoppingcart.CreateRequest{
		Items: shoppingcart.Items{
			shoppingcart.Item{
				Name:        testName,
				Price:       float32(testAmount),
				Description: testDescription,
			},
		},
	}

	cID, _ := uuid.NewUUID()
	c := coupon.Coupon{
		ID:     cID,
		Amount: 500,
	}
	createdShoppingCart := shoppingcart.New(req)

	t.Run("coupon already applied", func(t *testing.T) {
		sc := shoppingcart.ShoppingCart{
			CouponID: cID,
		}
		err := sc.ApplyCoupon(c.ID, c.Amount)
		assert.Equal(t, shoppingcart.ErrShoppinCartCouponAlreadyApplied, err)
	})

	t.Run("coupon excess the shopping cart amount", func(t *testing.T) {
		err := createdShoppingCart.ApplyCoupon(c.ID, c.Amount)
		assert.Equal(t, shoppingcart.ErrShoppointCartCouponAmountExceeded, err)
	})
}
