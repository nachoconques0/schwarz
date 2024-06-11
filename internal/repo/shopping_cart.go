package repo

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	internalErrors "github.com/nachoconques0/schwarz-challenge/internal/errors"
	shoppingcart "github.com/nachoconques0/schwarz-challenge/internal/shopping_cart"
)

var (
	// ErrMissingDB used when DB is nil
	ErrMissingDB = internalErrors.NewNotFound("DB connection is missing")
	// ErrShoppingCartsNotFound used when there are no shopping carts
	ErrShoppingCartsNotFound = internalErrors.NewNotFound("shopping carts not found")
	// ErrShoppingCartNotFound used when there is no shopping cart
	ErrShoppingCartNotFound = internalErrors.NewNotFound("shopping cart not found")
)

// shoppingCartTable is the table name for the shopping cart model
const shoppingCartTable = "schwarz.shopping_cart"

type shoppingRepository struct {
	db *gorm.DB
}

func NewShoppingCarRepository(db *gorm.DB) (shoppingcart.Repository, error) {
	if db == nil {
		return nil, ErrMissingDB
	}
	return &shoppingRepository{
		db: db,
	}, nil
}

// CreateShoppingCart will create a new shopping cart
func (sc shoppingRepository) CreateShoppingCart(shoppingCart *shoppingcart.ShoppingCart) (*shoppingcart.ShoppingCart, error) {
	if err := sc.db.Table(shoppingCartTable).Create(&shoppingCart).Error; err != nil {
		return nil, err
	}
	return shoppingCart, nil
}

// ListShoppingCarts returns a list of shopping carts
func (sc shoppingRepository) ListShoppingCarts() ([]shoppingcart.ShoppingCart, error) {
	var result []shoppingcart.ShoppingCart
	if err := sc.db.Table(shoppingCartTable).Find(&result).Error; err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrShoppingCartsNotFound
	}
	return result, nil
}

// GetShoppingCartForUpdate returns a shopping cart and it will lock the row in order to update it
func (sc shoppingRepository) GetShoppingCartForUpdate(tx *gorm.DB, scID uuid.UUID) (*shoppingcart.ShoppingCart, error) {
	var result shoppingcart.ShoppingCart

	if err := tx.Table(shoppingCartTable).
		Where(shoppingcart.ShoppingCart{ID: scID}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShoppingCartNotFound
		}
		return nil, err
	}

	return &result, nil
}

// UpdateShoppingCart updates shopping cart entity
func (sc shoppingRepository) UpdateShoppingCart(tx *gorm.DB, shoppingCart *shoppingcart.ShoppingCart) (*shoppingcart.ShoppingCart, error) {
	if err := tx.Table(shoppingCartTable).Save(&shoppingCart).Error; err != nil {
		return nil, err
	}
	return shoppingCart, nil
}

func (sc *shoppingRepository) BeginTransaction() *gorm.DB {
	return sc.db.Begin()
}

func (sc *shoppingRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (sc *shoppingRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}
