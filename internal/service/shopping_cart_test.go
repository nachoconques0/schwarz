package service_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/mocks"
	"github.com/nachoconques0/schwarz-challenge/internal/service"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
)

type testShoppingCartService struct {
	svc                  shoppingcart.Service
	couponMockRepo       *mocks.MockCouponRepository
	shoppingCartMockRepo *mocks.MockShoppingCartRepository
}

func TestShoppingCartService_CreateCoupon(t *testing.T) {
	ts := buildShoppingCartService(t)
	shoppingCartReq := &shoppingcart.CreateRequest{
		Items: shoppingcart.Items{
			shoppingcart.Item{
				Price:       10,
				Name:        "test",
				Description: "description",
			},
		},
	}
	expectedShoppingCart := shoppingcart.New(*shoppingCartReq)
	testCases := map[string]struct {
		req                  *shoppingcart.CreateRequest
		mocks                func()
		expectedShoppingCart *shoppingcart.ShoppingCart
		expectedError        error
	}{
		"repo fail": {
			req: shoppingCartReq,
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().CreateShoppingCart(gomock.Any()).Return(nil, errGeneric)
			},
			expectedShoppingCart: nil,
			expectedError:        errGeneric,
		},
		"invalid data": {
			req: &shoppingcart.CreateRequest{
				Items: shoppingcart.Items{},
			},
			mocks:                func() {},
			expectedShoppingCart: nil,
			expectedError:        shoppingcart.ErrShoppingCartEmptyItems,
		},
		"success": {
			req: shoppingCartReq,
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().CreateShoppingCart(gomock.Any()).Return(expectedShoppingCart, nil)
			},
			expectedShoppingCart: expectedShoppingCart,
			expectedError:        nil,
		},
	}

	for name, tc := range testCases {
		tc.mocks()
		t.Run(name, func(t *testing.T) {
			c, err := ts.svc.CreateShoppingCart(*tc.req)
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Equal(t, tc.expectedShoppingCart.Amount, c.Amount)
			}
		})
	}
}

func TestShoppingCartService_ListShoppingCarts(t *testing.T) {
	ts := buildShoppingCartService(t)

	createdSc := shoppingcart.New(shoppingcart.CreateRequest{
		Items: shoppingcart.Items{
			shoppingcart.Item{
				Price:       10,
				Name:        "test",
				Description: "description",
			},
		},
	})
	testCases := map[string]struct {
		mocks                 func()
		expectedShoppingCarts []shoppingcart.ShoppingCart
		expectedError         error
	}{
		"repo fail": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().ListShoppingCarts().Return(nil, errGeneric)
			},
			expectedShoppingCarts: []shoppingcart.ShoppingCart{},
			expectedError:         errGeneric,
		},
		"success": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().ListShoppingCarts().Return([]shoppingcart.ShoppingCart{*createdSc}, nil)
			},
			expectedShoppingCarts: []shoppingcart.ShoppingCart{
				*createdSc,
			},
			expectedError: nil,
		},
	}

	for name, tc := range testCases {
		tc.mocks()
		t.Run(name, func(t *testing.T) {
			res, err := ts.svc.ListShoppingCarts()
			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.Len(t, res, 1)
			}
		})
	}
}

func TestShoppingCartService_ApplyCoupon(t *testing.T) {
	ts := buildShoppingCartService(t)
	couponID := uuid.MustParse(uuid.NewString())
	scID := uuid.MustParse(uuid.NewString())

	invalidCoupon := &coupon.Coupon{
		Amount: 50,
		Used:   true,
	}

	toUpdateShoppingCart := &shoppingcart.ShoppingCart{
		CouponID: uuid.MustParse(uuid.Nil.String()),
		Items: shoppingcart.Items{
			shoppingcart.Item{
				Price:       100,
				Name:        "test",
				Description: "description",
			},
		},
		Amount: 100,
		Total:  100,
	}

	testCases := map[string]struct {
		req           *shoppingcart.CreateRequest
		mocks         func()
		expectedError error
	}{
		"GetCouponForUpdate fails": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().BeginTransaction().Return(nil)
				ts.couponMockRepo.EXPECT().GetCouponForUpdate(gomock.Any(), gomock.Any()).Return(nil, errGeneric)
			},
			expectedError: errGeneric,
		},
		"coupon is used": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().BeginTransaction().Return(nil)
				ts.couponMockRepo.EXPECT().GetCouponForUpdate(gomock.Any(), gomock.Any()).Return(invalidCoupon, nil)
			},
			expectedError: coupon.ErrCouponAlreadyUsed,
		},
		"GetShoppingCartForUpdate fails": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().BeginTransaction().Return(nil)
				ts.couponMockRepo.EXPECT().GetCouponForUpdate(gomock.Any(), gomock.Any()).Return(&coupon.Coupon{
					Amount: 50,
					Used:   false,
				}, nil)
				ts.shoppingCartMockRepo.EXPECT().GetShoppingCartForUpdate(gomock.Any(), gomock.Any()).Return(nil, errGeneric)
			},
			expectedError: errGeneric,
		},
		"shopping cart coupon has been applied": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().BeginTransaction().Return(nil)
				ts.couponMockRepo.EXPECT().GetCouponForUpdate(gomock.Any(), gomock.Any()).Return(&coupon.Coupon{
					Amount: 50,
					Used:   false,
				}, nil)
				ts.shoppingCartMockRepo.EXPECT().GetShoppingCartForUpdate(gomock.Any(), gomock.Any()).Return(&shoppingcart.ShoppingCart{
					CouponID: uuid.MustParse(uuid.NewString()),
					Items: shoppingcart.Items{
						shoppingcart.Item{
							Price:       100,
							Name:        "test",
							Description: "description",
						},
					},
					Amount: 100,
					Total:  100,
				}, nil)
			},
			expectedError: shoppingcart.ErrShoppinCartCouponAlreadyApplied,
		},
		"shopping cart update fails": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().BeginTransaction().Return(nil)
				ts.couponMockRepo.EXPECT().GetCouponForUpdate(gomock.Any(), gomock.Any()).Return(&coupon.Coupon{
					Amount: 50,
					Used:   false,
				}, nil)
				ts.shoppingCartMockRepo.EXPECT().GetShoppingCartForUpdate(gomock.Any(), gomock.Any()).Return(toUpdateShoppingCart, nil)
				ts.shoppingCartMockRepo.EXPECT().UpdateShoppingCart(gomock.Any(), gomock.Any()).Return(nil, errGeneric)
			},
			expectedError: errGeneric,
		},
		"success": {
			mocks: func() {
				ts.shoppingCartMockRepo.EXPECT().BeginTransaction().Return(nil)
				ts.couponMockRepo.EXPECT().GetCouponForUpdate(gomock.Any(), gomock.Any()).Return(&coupon.Coupon{
					Amount: 50,
					Used:   false,
				}, nil)
				ts.shoppingCartMockRepo.EXPECT().GetShoppingCartForUpdate(gomock.Any(), gomock.Any()).Return(&shoppingcart.ShoppingCart{
					CouponID: uuid.MustParse(uuid.Nil.String()),
					Items: shoppingcart.Items{
						shoppingcart.Item{
							Price:       100,
							Name:        "test",
							Description: "description",
						},
					},
					Amount: 100,
					Total:  100,
				}, nil)
				ts.shoppingCartMockRepo.EXPECT().UpdateShoppingCart(gomock.Any(), gomock.Any()).Return(toUpdateShoppingCart, nil)
				ts.couponMockRepo.EXPECT().UpdateCoupon(gomock.Any(), gomock.Any()).Return(nil, nil)
				ts.shoppingCartMockRepo.EXPECT().CommitTransaction(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
	}

	for name, tc := range testCases {
		tc.mocks()
		t.Run(name, func(t *testing.T) {
			err := ts.svc.ApplyCoupon(couponID, scID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func buildShoppingCartService(t *testing.T) testShoppingCartService {
	ctrl := gomock.NewController(t)
	couponRepo := mocks.NewMockCouponRepository(ctrl)
	shoppingCartRepo := mocks.NewMockShoppingCartRepository(ctrl)
	svc, err := service.NewShoppingCartService(shoppingCartRepo, couponRepo)
	assert.Nil(t, err)
	return testShoppingCartService{
		svc:                  svc,
		couponMockRepo:       couponRepo,
		shoppingCartMockRepo: shoppingCartRepo,
	}
}
