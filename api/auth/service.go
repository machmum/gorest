package auth

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/postgresql"
	"github.com/machmum/gorest/utl/platform/rds"
)

// Auth represents auth application service
type Auth struct {
	conn     Connection
	platform Platform
	cfg      *config.RDB
}

// Represent connection to used
type Connection struct {
	rds        *redis.Client
	postgresql *gorm.DB
}

// Represent database handler
type Platform struct {
	rds        RDSRefresh
	postgresql DBProfile
}

// Initialize initializes oauth application service
func Initialize(db *gorm.DB, rdb *redis.Client, cfgrdb *config.RDB) *Auth {
	return New(db, rdb, cfgrdb, postgresql.NewProfile(), rds.NewRefresh())
}

// New creates new iam service
func New(db *gorm.DB, rdb *redis.Client, cfg *config.RDB, postgresql DBProfile, rds RDSRefresh) *Auth {
	return &Auth{
		conn: Connection{
			rds:        rdb,
			postgresql: db,
		},
		platform: Platform{
			rds:        rds,
			postgresql: postgresql,
		},
		cfg: cfg,
	}
}

// Service represents auth service interface
type Service interface {
	Login(context echo.Context, username string, password string) (login gorest.Profile, err error)
	Logout(context echo.Context, token string) error
}

// DBProfile represent profile repository (postgresql) interface
type DBProfile interface {
	FindByUsername(postgresql *gorm.DB, username string) (result gorestdb.Profile, err error)
}

// RDSRefresh represent refresh repository (redis) interface
type RDSRefresh interface {
	Save(rdb *redis.Client, data interface{}, key, index string, lifetime uint) error
	Delete(rdb *redis.Client, key, index string) error
}
