package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	internalHTTP "github.com/nachoconques0/schwarz-challenge/internal/http"
	"github.com/nachoconques0/schwarz-challenge/internal/mocks"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestController_CreateShoppingCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockShoppingCartService(ctrl)
	controller := internalHTTP.NewShopppingCartCtrl(svc)

	t.Run("success", func(t *testing.T) {
		svc.EXPECT().CreateShoppingCart(gomock.Any()).Return(&shoppingcart.ShoppingCart{
			Items: shoppingcart.Items{
				shoppingcart.Item{
					Price: float32(testAmount),
					Name:  testName,
				},
			},
			Amount: float32(testAmount),
		}, nil)

		body, err := json.Marshal(shoppingcart.CreateRequest{
			Items: shoppingcart.Items{
				shoppingcart.Item{
					Price: float32(testAmount),
					Name:  testName,
				},
			},
		})
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", bytes.NewBuffer(body))
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.CreateShoppingCart(recorder, req)
		resp := recorder.Result()

		response := &shoppingcart.ShoppingCart{}
		err = json.NewDecoder(resp.Body).Decode(response)
		assert.Nil(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, response.Amount, float32(testAmount))
		_ = resp.Body.Close()
	})

	t.Run("fail", func(t *testing.T) {
		svc.EXPECT().CreateShoppingCart(gomock.Any()).Return(nil, errTest)
		body, err := json.Marshal(shoppingcart.CreateRequest{
			Items: shoppingcart.Items{
				shoppingcart.Item{
					Price: float32(testAmount),
					Name:  testName,
				},
			},
		})
		assert.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", bytes.NewBuffer(body))
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.CreateShoppingCart(recorder, req)
		resp := recorder.Result()
		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, errTest, responseErr)
		_ = resp.Body.Close()
	})
}

func TestController_ListShoppingCarts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockShoppingCartService(ctrl)
	controller := internalHTTP.NewShopppingCartCtrl(svc)

	t.Run("success", func(t *testing.T) {
		svc.EXPECT().ListShoppingCarts().Return([]shoppingcart.ShoppingCart{
			{
				Items: shoppingcart.Items{
					shoppingcart.Item{
						Price: float32(testAmount),
						Name:  testName,
					},
				},
				Amount: float32(testAmount),
			},
		}, nil)

		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", nil)
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.ListShoppingCarts(recorder, req)
		resp := recorder.Result()

		response := []shoppingcart.ShoppingCart{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.Nil(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, response[0].Amount, float32(testAmount))
		_ = resp.Body.Close()
	})

	t.Run("fail", func(t *testing.T) {
		svc.EXPECT().ListShoppingCarts().Return(nil, errTest)

		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", nil)
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.ListShoppingCarts(recorder, req)
		resp := recorder.Result()
		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, errTest, responseErr)
		_ = resp.Body.Close()
	})
}

func TestController_ApplyCoupon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	couponID, _ := uuid.NewUUID()
	shoppingCartID, _ := uuid.NewUUID()
	svc := mocks.NewMockShoppingCartService(ctrl)
	controller := internalHTTP.NewShopppingCartCtrl(svc)

	t.Run("success", func(t *testing.T) {
		svc.EXPECT().ApplyCoupon(shoppingCartID, couponID).Return(nil)
		urlVars := map[string]string{
			"id":        shoppingCartID.String(),
			"coupon_id": couponID.String(),
		}

		req, err := http.NewRequest(http.MethodGet, "http://www.test.com", nil)
		assert.Nil(t, err)
		req = mux.SetURLVars(req, urlVars)
		q := req.URL.Query()
		req.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()
		controller.ApplyCoupon(recorder, req)

		resp := recorder.Result()
		assert.Equal(t, recorder.Code, http.StatusOK)
		_ = resp.Body.Close()
	})

	t.Run("missing shopping cart id", func(t *testing.T) {
		urlVars := map[string]string{
			"id":        "",
			"coupon_id": couponID.String(),
		}

		req, err := http.NewRequest(http.MethodGet, "http://www.test.com", nil)
		assert.Nil(t, err)
		req = mux.SetURLVars(req, urlVars)
		q := req.URL.Query()
		req.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()
		controller.ApplyCoupon(recorder, req)

		resp := recorder.Result()
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, internalHTTP.ErrShoppingCartEmptyID, responseErr)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		_ = resp.Body.Close()
	})

	t.Run("missing coupon id", func(t *testing.T) {
		urlVars := map[string]string{
			"id": shoppingCartID.String(),
		}

		req, err := http.NewRequest(http.MethodGet, "http://www.test.com", nil)
		assert.Nil(t, err)
		req = mux.SetURLVars(req, urlVars)
		q := req.URL.Query()
		req.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()
		controller.ApplyCoupon(recorder, req)

		resp := recorder.Result()
		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, internalHTTP.ErrCouponEmptyID, responseErr)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		_ = resp.Body.Close()
	})

	t.Run("fail svc", func(t *testing.T) {
		svc.EXPECT().ApplyCoupon(shoppingCartID, couponID).Return(errTest)
		urlVars := map[string]string{
			"id":        shoppingCartID.String(),
			"coupon_id": couponID.String(),
		}

		req, err := http.NewRequest(http.MethodGet, "http://www.test.com", nil)
		assert.Nil(t, err)
		req = mux.SetURLVars(req, urlVars)
		q := req.URL.Query()
		req.URL.RawQuery = q.Encode()

		recorder := httptest.NewRecorder()
		controller.ApplyCoupon(recorder, req)

		resp := recorder.Result()
		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, errTest, responseErr)
		_ = resp.Body.Close()
	})

}
