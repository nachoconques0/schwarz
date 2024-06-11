package repo_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/helpers"
	"github.com/nachoconques0/schwarz-challenge/internal/repo"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	couponID   = uuid.New()
	testCoupon = &coupon.Coupon{
		ID:     couponID,
		Name:   "couponName",
		Amount: 10,
		Used:   false,
	}
)

func TestRepository_CreateCoupon(t *testing.T) {
	c := &coupon.Coupon{
		ID:     couponID,
		Name:   "couponName",
		Amount: 10,
		Used:   false,
	}

	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	r := createCouponRepo(t, db)

	t.Run("it should create the coupon", func(t *testing.T) {
		res, err := r.CreateCoupon(c)
		assert.Nil(t, err)
		assert.Equal(t, c.ID, res.ID)
		assert.Equal(t, c.Name, res.Name)
		assert.Equal(t, c.Amount, res.Amount)
		assert.Equal(t, c.Used, res.Used)
		assert.NotEqual(t, res.CreatedAt, time.Time{})
		assert.NotEqual(t, res.UpdatedAt, time.Time{})
	})

	t.Run("it should fail if the coupon already exists", func(t *testing.T) {
		updatedCoupon := c
		updatedCoupon.Name = "new name"
		_, err := r.CreateCoupon(updatedCoupon)
		assert.NotNil(t, err)
	})
}

func TestRepository_ListCoupon(t *testing.T) {
	testCases := map[string]struct {
		expectedError error
		expectedLen   int
	}{
		"when the are no coupons": {
			expectedError: repo.ErrCouponNotFound,
			expectedLen:   0,
		},
		"when coupon exists": {
			expectedError: nil,
			expectedLen:   1,
		},
	}

	for name, tc := range testCases {
		db, teardown, err := helpers.NewTestDB()
		if err != nil {
			assert.Nil(t, err)
		}
		defer teardown()

		r := createCouponRepo(t, db)
		if tc.expectedLen != 0 {
			createCoupon(t, r)
		}

		t.Run(name, func(t *testing.T) {
			res, err := r.ListCoupons()
			assert.Equal(t, tc.expectedError, err)
			if res != nil {
				assert.Len(t, res, tc.expectedLen)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func TestRepository_GetCouponForUpdate(t *testing.T) {
	db, teardown, err := helpers.NewTestDB()
	if err != nil {
		assert.Nil(t, err)
	}
	defer teardown()

	r := createCouponRepo(t, db)

	createdCoupon := createCoupon(t, r)

	testCases := map[string]struct {
		expectedError  error
		expectedCoupon *coupon.Coupon
		id             uuid.UUID
	}{
		"when there is no coupon": {
			id:             uuid.New(),
			expectedError:  repo.ErrCouponNotFound,
			expectedCoupon: nil,
		},
		"when coupon exists": {
			id:             createdCoupon.ID,
			expectedError:  nil,
			expectedCoupon: testCoupon,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res, err := r.GetCouponForUpdate(db, tc.id)
			assert.Equal(t, tc.expectedError, err)
			if res != nil {
				assert.Equal(t, tc.expectedCoupon.ID, res.ID)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func TestRepository_UpdateCoupon(t *testing.T) {
	testCases := map[string]struct {
		expectedError error
		isUsed        bool
	}{
		"when coupon exists": {
			expectedError: nil,
			isUsed:        true,
		},
	}
	for name, tc := range testCases {
		db, teardown, err := helpers.NewTestDB()
		if err != nil {
			assert.Nil(t, err)
		}
		defer teardown()

		r := createCouponRepo(t, db)
		createdCoupon := createCoupon(t, r)

		t.Run(name, func(t *testing.T) {
			createdCoupon.Used = true
			res, err := r.UpdateCoupon(db, createdCoupon)
			assert.Equal(t, tc.expectedError, err)
			if res != nil {
				assert.True(t, res.Used)
			} else {
				assert.Nil(t, res)
			}
		})
	}
}

func createCouponRepo(t *testing.T, db *gorm.DB) coupon.Repository {
	r, err := repo.NewCouponRepository(db)
	if err != nil {
		assert.Nil(t, err)
	}
	return r
}

func createCoupon(t *testing.T, r coupon.Repository) *coupon.Coupon {
	c, err := r.CreateCoupon(testCoupon)
	if err != nil {
		assert.Nil(t, err)
	}
	return c
}
