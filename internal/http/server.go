// Package http provides a http server implementation for the backoffice-svc.
// It defines all the routes available and the permissions needed by the user to access them.
package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nachoconques0/schwarz-challenge/internal/coupon"
	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
)

// Server type holds the dependencies needed
// for handle an http.Server
type Server struct {
	*http.Server
	shoppingCartSrv shoppingcart.Server
	couponSrv       coupon.Server
}

// NewServer builds a new http.Server by using the given dependencies
// all of thoses dependencies are mandatory
func NewServer(
	port string,
	shoppingCartSrv shoppingcart.Server,
	couponSrv coupon.Server,
) (*Server, error) {
	if port == "" {
		return nil, errors.New("server port can not be empty")
	}
	if shoppingCartSrv == nil {
		return nil, errors.New("shopping cart server can not be nil")
	}
	if couponSrv == nil {
		return nil, errors.New("coupon server can not be nil")
	}
	s := &Server{
		shoppingCartSrv: shoppingCartSrv,
		couponSrv:       couponSrv,
	}
	s.Server = &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
	}
	s.router()

	return s, nil
}

// Run method starts the http.Server, so it's
// ready for receive http.Request
func (s *Server) Run() error {
	s.Handler = handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)(s.Handler)
	if err := s.ListenAndServe(); err != nil {
		return fmt.Errorf("error while starting the http server %s", err)
	}
	return nil
}

// Stop method will shutdown the http.Server, after
// this method is called we are no longer available
// for handle more http.Request
func (s *Server) Stop(ctx context.Context) error {
	if err := s.Shutdown(ctx); err != nil {
		return fmt.Errorf("error while shutting down the http server %s", err)
	}
	return nil
}

func (s *Server) router() {
	r := mux.NewRouter()
	r.NewRoute().Subrouter()

	s.shoppingCartRouter(r)
	s.couponRouter(r)

	r.Use(contentTypeJSONMiddleware)
	// Pass our instance of gorilla/mux in.
	s.Handler = r
}

// shoppingCartRouter holds the routing for the shopping cart endpoints
func (s *Server) shoppingCartRouter(r *mux.Router) {
	r.HandleFunc("/shopping-cart", s.shoppingCartSrv.CreateShoppingCart).Methods(http.MethodPost)
	r.HandleFunc("/shopping-cart", s.shoppingCartSrv.ListShoppingCarts).Methods(http.MethodGet)
	r.HandleFunc("/shopping-cart/{id}/apply-coupon/{coupon_id}", s.shoppingCartSrv.ApplyCoupon).Methods(http.MethodPut)
}

// couponRouter holds the routing for the coupon endpoints
func (s *Server) couponRouter(r *mux.Router) {
	r.HandleFunc("/coupon", s.couponSrv.CreateCoupon).Methods(http.MethodPost)
	r.HandleFunc("/coupon", s.couponSrv.LisCoupons).Methods(http.MethodGet)
}

func contentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// encodeResponse receives the http response writer and the response
// to be encoded. It also sets the StatusCode to 200 unless encoding fails, in that
// case it encodes a code 400 and the error
func encodeResponse(w http.ResponseWriter, res interface{}) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		slog.Error(fmt.Sprintf("error encoding response: %s", err))
	}
}

// responseError handles internals error http response.
func responseError(w http.ResponseWriter, r *http.Request, err error) {
	var internalErr *internalErrors.Error
	if errors.As(err, &internalErr) {
		w.WriteHeader(internalErr.HTTPStatus())
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(internalErr); err != nil {
			slog.Error(fmt.Sprintf("error encoding response: %s", err))
		}
		return
	}
	responseError(w, r, err)
}
