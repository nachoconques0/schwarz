package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	internalHTTP "github.com/nachoconques0/schwarz-challenge/internal/http"
	"github.com/nachoconques0/schwarz-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	errTest    = errors.NewInternalError("expected error")
	testName   = "testName"
	testAmount = 100
)

func TestController_CreateCoupon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockCouponService(ctrl)
	controller := internalHTTP.NewCouponCtrl(svc)

	bodyParams := map[string]interface{}{
		"name":   testName,
		"amount": testAmount,
	}

	t.Run("success", func(t *testing.T) {
		svc.EXPECT().CreateCoupon(gomock.Any()).Return(&coupon.Coupon{
			Name:   testName,
			Amount: float32(testAmount),
		}, nil)

		body, _ := json.Marshal(bodyParams)
		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", bytes.NewBuffer(body))
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.CreateCoupon(recorder, req)
		resp := recorder.Result()

		response := &coupon.Coupon{}
		err = json.NewDecoder(resp.Body).Decode(response)
		assert.Nil(t, err)
		assert.Equal(t, response.Name, testName)
		assert.Equal(t, response.Amount, float32(testAmount))
		_ = resp.Body.Close()
	})

	t.Run("fail", func(t *testing.T) {
		svc.EXPECT().CreateCoupon(gomock.Any()).Return(nil, errTest)
		body, _ := json.Marshal(bodyParams)
		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", bytes.NewBuffer(body))
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.CreateCoupon(recorder, req)
		resp := recorder.Result()

		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, errTest, responseErr)
		_ = resp.Body.Close()
	})
}

func TestController_LisCoupons(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	svc := mocks.NewMockCouponService(ctrl)
	controller := internalHTTP.NewCouponCtrl(svc)

	t.Run("success", func(t *testing.T) {
		svc.EXPECT().ListCoupons().Return([]coupon.Coupon{
			{
				Name:   testName,
				Amount: float32(testAmount),
			},
		}, nil)

		req, err := http.NewRequest(http.MethodGet, "http://www.test.com", nil)
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.LisCoupons(recorder, req)
		resp := recorder.Result()

		response := []coupon.Coupon{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.Nil(t, err)
		assert.Len(t, response, 1)
		_ = resp.Body.Close()
	})

	t.Run("fail", func(t *testing.T) {
		svc.EXPECT().ListCoupons().Return(nil, errTest)
		req, err := http.NewRequest(http.MethodPost, "http://www.test.com", nil)
		assert.Nil(t, err)

		recorder := httptest.NewRecorder()
		controller.LisCoupons(recorder, req)
		resp := recorder.Result()

		responseErr := &errors.Error{}
		err = json.NewDecoder(resp.Body).Decode(responseErr)
		assert.Nil(t, err)
		assert.Equal(t, errTest, responseErr)
		_ = resp.Body.Close()
	})
}
