package service_test

import (
	"testing"

	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	"github.com/nachoconques0/schwarz-challenge/internal/mocks"
	"github.com/nachoconques0/schwarz-challenge/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	errGeneric = errors.NewInternalError("some-error")
)

type testCouponService struct {
	svc            coupon.Service
	couponMockRepo *mocks.MockCouponRepository
}

func TestCouponService_CreateCoupon(t *testing.T) {
	ts := buildCouponService(t)
	couponReq := &coupon.CreateRequest{
		Name:   "test",
		Amount: 100,
	}
	expectedCoupon := coupon.New(*couponReq)
	testCases := map[string]struct {
		req            *coupon.CreateRequest
		mocks          func()
		expectedCoupon *coupon.Coupon
		expectedError  error
	}{
		"repo fail": {
			req: couponReq,
			mocks: func() {
				ts.couponMockRepo.EXPECT().CreateCoupon(gomock.Any()).Return(nil, errGeneric)
			},
			expectedCoupon: nil,
			expectedError:  errGeneric,
		},
		"invalid data": {
			req: &coupon.CreateRequest{
				Name:   "test",
				Amount: 0,
			},
			mocks:          func() {},
			expectedCoupon: nil,
			expectedError:  coupon.ErrCouponInvalidAmount,
		},
		"success": {
			req: couponReq,
			mocks: func() {
				ts.couponMockRepo.EXPECT().CreateCoupon(gomock.Any()).Return(expectedCoupon, nil)
			},
			expectedCoupon: expectedCoupon,
			expectedError:  nil,
		},
	}

	for name, tc := range testCases {
		tc.mocks()
		t.Run(name, func(t *testing.T) {
			c, err := ts.svc.CreateCoupon(*tc.req)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedCoupon.Name, c.Name)
				assert.Equal(t, tc.expectedCoupon.Amount, c.Amount)
				assert.False(t, c.Used)
			}
		})
	}
}

func TestCouponService_ListCoupons(t *testing.T) {
	ts := buildCouponService(t)
	couponReq := &coupon.CreateRequest{
		Name:   "test",
		Amount: 100,
	}
	c := coupon.New(*couponReq)
	testCases := map[string]struct {
		mocks           func()
		expectedCoupons []coupon.Coupon
		expectedError   error
	}{
		"repo fail": {
			mocks: func() {
				ts.couponMockRepo.EXPECT().ListCoupons().Return(nil, errGeneric)
			},
			expectedCoupons: []coupon.Coupon{},
			expectedError:   errGeneric,
		},
		"success": {
			mocks: func() {
				ts.couponMockRepo.EXPECT().ListCoupons().Return([]coupon.Coupon{*c}, nil)
			},
			expectedCoupons: []coupon.Coupon{
				*c,
			},
			expectedError: nil,
		},
	}

	for name, tc := range testCases {
		tc.mocks()
		t.Run(name, func(t *testing.T) {
			res, err := ts.svc.ListCoupons()
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Len(t, res, 1)
			}
		})
	}
}

func buildCouponService(t *testing.T) testCouponService {
	ctrl := gomock.NewController(t)
	couponRepo := mocks.NewMockCouponRepository(ctrl)
	svc, err := service.NewCouponService(couponRepo)
	assert.Nil(t, err)

	return testCouponService{
		svc:            svc,
		couponMockRepo: couponRepo,
	}
}
