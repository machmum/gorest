package product

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/postgresql"
)

// Auth represents auth application service
type Product struct {
	conn     Connection
	platform Platform
	cfg      *config.RDB
}

type Connection struct {
	postgresql *gorm.DB
	rds        *redis.Client
}

type Platform struct {
	Product DBProduct
}

// New creates new iam service
func New(db *gorm.DB, rdb *redis.Client, cfg *config.RDB, postgresql *Platform) *Product {
	return &Product{
		conn: Connection{
			postgresql: db,
			rds:        rdb,
		},
		platform: Platform{
			Product: postgresql.Product,
		},
		cfg: cfg,
	}
}

func NewPostgresql() *Platform {
	return &Platform{
		Product: postgresql.NewProduct(),
	}
}

// Initialize initializes oauth application service
func Initialize(db *gorm.DB, rdb *redis.Client, cfgrdb *config.RDB) *Product {
	return New(db, rdb, cfgrdb, NewPostgresql())
}

// Service represents auth service interface
type Service interface {
	Product(c echo.Context, profileID int, pid int) (products gorest.Product, err error)
}

type DBProduct interface {
	FindByProductID(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error)
}
