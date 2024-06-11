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

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"

	"github.com/nachoconques0/schwarz-challenge/internal/errors"
	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
)

var (
	// ErrShoppingCartEmptyID used when shopping cart ID is invalid
	ErrShoppingCartEmptyID = internalErrors.NewWrongInput("shopping cart ID is invalid")
	// ErrCouponEmptyID used when coupon ID is invalid
	ErrCouponEmptyID = internalErrors.NewWrongInput("coupon ID is invalid")
	// ErrInvalidCreateShoppingCartRequest used when create shopping cart request contains invalid data
	ErrInvalidCreateShoppingCartRequest = errors.NewWrongInput("invalid create shopping cart request")
)

//go:embed schemas/shopping_cart/create.json
var createShoppingCartRequestSchema []byte

// NewShopppingCartCtrl creates a new HTTP Controller
// with the given shoppingcart.Service
func NewShopppingCartCtrl(svc shoppingcart.Service) shoppingcart.Server {
	return &shoppingCartController{svc: svc}
}

// shoppingCartController holds the required dependencies
// in order to implement the service Request
type shoppingCartController struct {
	svc shoppingcart.Service
}

// CreateShoppingCart receives a request in order to create a shopping cart
func (scCtrl *shoppingCartController) CreateShoppingCart(w http.ResponseWriter, r *http.Request) {
	createSchema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(createShoppingCartRequestSchema))
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
		responseError(w, r, ErrInvalidCreateShoppingCartRequest)
		return
	}

	var payload shoppingcart.CreateRequest
	err = json.Unmarshal(requestBytes, &payload)
	if err != nil {
		slog.Error(fmt.Sprintf("decoding create shopping cart request: %s\n", err))
		responseError(w, r, err)
		return
	}

	res, err := scCtrl.svc.CreateShoppingCart(payload)
	if err != nil {
		slog.Error(fmt.Sprintf("creating shopping cart: %s\n", err))
		responseError(w, r, err)
		return
	}

	encodeResponse(w, res)
}

// ListShoppingCarts returns a list of shopping carts
func (scCtrl *shoppingCartController) ListShoppingCarts(w http.ResponseWriter, r *http.Request) {
	res, err := scCtrl.svc.ListShoppingCarts()
	if err != nil {
		slog.Error(fmt.Sprintf("listing shopping cart: %s\n", err))
		responseError(w, r, err)
		return
	}
	encodeResponse(w, res)
}

// ApplyCoupon receives a request in order to apply a coupon to a shopping cart
func (scCtrl *shoppingCartController) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shoppingCartID := vars["id"]
	couponID := vars["coupon_id"]

	if shoppingCartID == "" {
		slog.Error(fmt.Sprintf("ctrl: applying coupon: %s\n", ErrShoppingCartEmptyID))
		responseError(w, r, ErrShoppingCartEmptyID)
		return
	}
	if couponID == "" {
		slog.Error(fmt.Sprintf("ctrl: applying coupon: %s\n", ErrCouponEmptyID))
		responseError(w, r, ErrCouponEmptyID)
		return
	}

	parsedShoppingCartID, err := uuid.Parse(shoppingCartID)
	if err != nil {
		slog.Error(fmt.Sprintf("ctrl: applying coupon: %s\n", ErrShoppingCartEmptyID))
		w.WriteHeader(http.StatusBadRequest)
		encodeResponse(w, nil)
		return
	}
	parsedCouponID, err := uuid.Parse(couponID)
	if err != nil {
		slog.Error(fmt.Sprintf("ctrl: applying coupon: %s\n", ErrCouponEmptyID))
		w.WriteHeader(http.StatusBadRequest)
		encodeResponse(w, nil)
		return
	}
	err = scCtrl.svc.ApplyCoupon(parsedShoppingCartID, parsedCouponID)
	if err != nil {
		slog.Error(fmt.Sprintf("ctrl: applying coupon: %s\n", err))
		responseError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
