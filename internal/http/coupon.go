package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	// embed used for loading request cases
	_ "embed"

	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	"github.com/xeipuuv/gojsonschema"
)

// ErrInvalidCreateCouponRequest used when create coupon request contains invalid data
var ErrInvalidCreateCouponRequest = errors.NewWrongInput("invalid create coupon request")

//go:embed schemas/coupon/create.json
var createRequestSchema []byte

// NewCouponCtrl creates a new HTTP Controller
// with the given coupon.Service
func NewCouponCtrl(svc coupon.Service) coupon.Server {
	return &couponController{svc: svc}
}

// couponController holds the required dependencies
// in order to implement the service Request
type couponController struct {
	svc coupon.Service
}

// CreateCoupon receives a request in order to create a coupon
func (cCtrl *couponController) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	createSchema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(createRequestSchema))
	if err != nil {
		slog.Error(fmt.Sprintf("creating create coupon schema : %s\n", err))
		responseError(w, r, err)
		return
	}
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error(fmt.Sprintf("reading create coupon body : %s\n", err))
		responseError(w, r, err)
		return
	}

	requestJSON := gojsonschema.NewBytesLoader(requestBytes)
	result, err := createSchema.Validate(requestJSON)
	if err != nil {
		slog.Error(fmt.Sprintf("error validating request : %s\n", err))
		responseError(w, r, err)
		return
	}

	if !result.Valid() {
		details := make([]string, 0, len(result.Errors()))
		for _, err := range result.Errors() {
			details = append(details, fmt.Sprintf("Field:%s, with error:%s:", err.Field(), err.Description()))
		}
		slog.Error("create coupon request data is not valid:",
			slog.String("error_details", strings.Join(details, "\n")),
		)
		responseError(w, r, ErrInvalidCreateCouponRequest)
		return
	}

	var payload coupon.CreateRequest
	err = json.Unmarshal(requestBytes, &payload)
	if err != nil {
		slog.Error(fmt.Sprintf("decoding create coupon request: %s\n", err))
		responseError(w, r, err)
		return
	}
	res, err := cCtrl.svc.CreateCoupon(payload)
	if err != nil {
		slog.Error(fmt.Sprintf("creating coupon cart: %s\n", err))
		responseError(w, r, err)
		return
	}

	encodeResponse(w, res)
}

// LisCoupons receives a request in order to list coupons
func (cCtrl *couponController) LisCoupons(w http.ResponseWriter, r *http.Request) {
	res, err := cCtrl.svc.ListCoupons()
	if err != nil {
		slog.Error(fmt.Sprintf("listing coupons: %s\n", err))
		responseError(w, r, err)
		return
	}
	encodeResponse(w, res)
}
