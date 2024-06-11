package app

import (
	"fmt"

	"github.com/nachoconques0/schwarz-challenge/internal/http"
)

// setupHTTPServer creates the server that handles all the
// http requests received
func (a *Application) setupHTTPServer() (err error) {
	scCtrl := http.NewShopppingCartCtrl(a.shoppingCartService)
	cCtrl := http.NewCouponCtrl(a.couponService)

	// a.server = &wrf.Server{}
	res, err := http.NewServer(a.httpPort, scCtrl, cCtrl)
	if err != nil {
		return fmt.Errorf("app: error setting up the http server %s", err)
	}
	a.server = res
	return nil
}
