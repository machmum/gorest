package profile

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/product"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/postgresql"
)

// Auth represents auth application service
type Profile struct {
	conn     Connection
	platform Platform
	cfg      *config.RDB
	internal Internal
}

type Connection struct {
	postgresql *gorm.DB
	rds        *redis.Client
}

type Platform struct {
	Profile DBProfile
	// product DBProduct
}

type Internal struct {
	Product ServiceProduct
}

// New creates new iam service
func New(db *gorm.DB, rdb *redis.Client, cfg *config.RDB, postgresql *Platform, internal *Internal) *Profile {
	return &Profile{
		conn: Connection{
			postgresql: db,
			rds:        rdb,
		},
		platform: Platform{
			Profile: postgresql.Profile,
			// product: postgresql.product,
		},
		cfg: cfg,
		internal: Internal{
			Product: internal.Product,
		},
	}
}

func NewPostgresql() *Platform {
	return &Platform{
		Profile: postgresql.NewProfile(),
		// product: postgresql.NewProduct(),
	}
}

func NewInternal(db *gorm.DB, rdb *redis.Client, cfg *config.RDB) *Internal {
	return &Internal{
		Product: product.Initialize(db, rdb, cfg),
	}
}

// Initialize initializes oauth application service
func Initialize(db *gorm.DB, rdb *redis.Client, cfgrdb *config.RDB) *Profile {
	return New(db, rdb, cfgrdb, NewPostgresql(), NewInternal(db, rdb, cfgrdb))
}

// Service represents auth service interface
type Service interface {
	Profile(ctx echo.Context, cturl string) (result gorest.ProfileSlice, err error)
	ProfileSimple(c echo.Context, profileID int) (profile gorest.Profile, err error)
	Product(c echo.Context, profileID int, pid int) (products gorest.Product, err error)
}

// Service product implement get product from service product
type ServiceProduct interface {
	Product(c echo.Context, profileID int, pid int) (products gorest.Product, err error)
}

// DBProfile represent profile repository (profile) in DB interface
type DBProfile interface {
	// db profile
	FindByCategoryReturnSimple(postgresql *gorm.DB, cturl string) (result gorestdb.ProfileSimpleSlice, err error)
	FindByProfileIDReturnSimple(postgresql *gorm.DB, id int) (result gorestdb.ProfileSimple, err error)
}

type DBProduct interface {
	FindByProductID(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error)
}
