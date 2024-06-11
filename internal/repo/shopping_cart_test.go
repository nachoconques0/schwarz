package repo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nachoconques0/schwarz-challenge/internal/helpers"
	"github.com/nachoconques0/schwarz-challenge/internal/repo"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	shoppingCartID   = uuid.New()
	testShoppingCart = &shoppingcart.ShoppingCart{
		ID:     shoppingCartID,
		Amount: 10,
		Items: shoppingcart.Items{
			shoppingcart.Item{Price: 10},
		},
	}
)

func TestRepository_CreateShoppingCart(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	r := createShoppingCartRepo(t, db)

	t.Run("it should create the shopping cart", func(t *testing.T) {
		res, err := r.CreateShoppingCart(testShoppingCart)
		assert.Nil(t, err)
		assert.Equal(t, testShoppingCart.ID, res.ID)
		assert.Equal(t, testShoppingCart.Amount, res.Amount)
		assert.Equal(t, testShoppingCart.Items, res.Items)
		assert.NotEqual(t, res.CreatedAt, time.Time{})
		assert.NotEqual(t, res.UpdatedAt, time.Time{})
	})

	t.Run("it should fail if the shopping cart already exists", func(t *testing.T) {
		updatedShoppingCart := testShoppingCart
		updatedShoppingCart.Amount = 30
		_, err := r.CreateShoppingCart(updatedShoppingCart)
		assert.NotNil(t, err)
	})
}

func TestRepository_ListShoppingCarts(t *testing.T) {
	testCases := map[string]struct {
		expectedError error
		expectedLen   int
	}{
		"when the are no shopping carts": {
			expectedError: repo.ErrShoppingCartsNotFound,
			expectedLen:   0,
		},
		"when shopping cart exists": {
			expectedError: nil,
			expectedLen:   1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			db, teardown, err := helpers.NewTestDB()
			if err != nil {
				assert.Nil(t, err)
			}
			defer teardown()

			r := createShoppingCartRepo(t, db)

			if tc.expectedLen != 0 {
				createShoppingCart(t, r)
			}
			res, err := r.ListShoppingCarts()
			assert.Equal(t, tc.expectedError, err)
			if res != nil {
				assert.Len(t, res, tc.expectedLen)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func TestRepository_GetShoppingCartForUpdate(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	r := createShoppingCartRepo(t, db)
	createdShoppingCart := createShoppingCart(t, r)

	testCases := map[string]struct {
		expectedError        error
		expectedShoppingCart *shoppingcart.ShoppingCart
		id                   uuid.UUID
	}{
		"when there is no shopping cart": {
			expectedError:        repo.ErrShoppingCartNotFound,
			expectedShoppingCart: nil,
			id:                   uuid.New(),
		},
		"when shopping cart exists": {
			expectedError:        nil,
			expectedShoppingCart: testShoppingCart,
			id:                   createdShoppingCart.ID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res, err := r.GetShoppingCartForUpdate(db, tc.id)
			assert.Equal(t, tc.expectedError, err)
			if res != nil {
				assert.Equal(t, tc.expectedShoppingCart.ID, res.ID)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func TestRepository_UpdateShoppingCart(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	r := createShoppingCartRepo(t, db)

	createdShoppingCart := createShoppingCart(t, r)
	testCases := map[string]struct {
		expectedError  error
		expectedAmount float32
	}{
		"when shopping cart exists": {
			expectedError:  nil,
			expectedAmount: 100,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			createdShoppingCart.Amount = 100
			res, err := r.UpdateShoppingCart(db, createdShoppingCart)
			assert.Equal(t, tc.expectedError, err)
			if res != nil {
				assert.Equal(t, tc.expectedAmount, res.Amount)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func createShoppingCartRepo(t *testing.T, db *gorm.DB) shoppingcart.Repository {
	r, err := repo.NewShoppingCarRepository(db)
	if err != nil {
		assert.Nil(t, err)
	}
	return r
}

func createShoppingCart(t *testing.T, r shoppingcart.Repository) *shoppingcart.ShoppingCart {
	sc, err := r.CreateShoppingCart(testShoppingCart)
	if err != nil {
		assert.Nil(t, err)
	}
	return sc
}
