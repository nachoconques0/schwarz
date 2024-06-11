package app

import (
	"github.com/nachoconques0/schwarz-challenge/internal/repo"
	"github.com/nachoconques0/schwarz-challenge/internal/service"
	"gorm.io/gorm"
)

// setupDomain will take the db as a given parameters, and with it will
// start all the dependency injections required for start up our business domain, the outcome
// will be a ready domain service or an error
func (a *Application) setupDomain(db *gorm.DB) error {
	couponRepo, err := repo.NewCouponRepository(db)
	if err != nil {
		return err
	}
	a.couponRepo = couponRepo

	shoppingCartRepoRepo, err := repo.NewShoppingCarRepository(db)
	if err != nil {
		return err
	}
	a.shoppingCartRepo = shoppingCartRepoRepo

	couponSvc, err := service.NewCouponService(a.couponRepo)
	if err != nil {
		return err
	}
	a.couponService = couponSvc

	scService, err := service.NewShoppingCartService(
		a.shoppingCartRepo,
		a.couponRepo,
	)
	if err != nil {
		return err
	}
	a.shoppingCartService = scService

	return nil
}
