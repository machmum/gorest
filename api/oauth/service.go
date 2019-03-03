package oauth

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/postgresql"
	"github.com/machmum/gorest/utl/platform/rds"
	"github.com/machmum/gorest/utl/secure"
)

// Auth represents auth application service
type Oauth struct {
	conn     Connection
	cfg      *config.RDB
	platform Platform
	token    TokenGenerator
}

type Connection struct {
	postgresql *gorm.DB
	rds        *redis.Client
}

type Platform struct {
	postgresql DBChannel
	rds        RDSRefresh
}

// New creates new iam service
func New(db *gorm.DB, rdb *redis.Client, cfg *config.RDB, postgresql DBChannel, rds RDSRefresh, tg TokenGenerator) *Oauth {
	return &Oauth{
		conn: Connection{
			postgresql: db,
			rds:        rdb,
		},
		cfg: cfg,
		platform: Platform{
			postgresql: postgresql,
			rds:        rds,
		},
		token: tg,
	}
}

// Initialize initializes oauth application service
func Initialize(db *gorm.DB, rdb *redis.Client, cfg *config.RDB) *Oauth {
	return New(db, rdb, cfg, postgresql.NewChannel(), rds.NewRefresh(), secure.NewToken())
}

// Service represents auth service interface
type Service interface {
	Tokenize(echo.Context, string, string, string, string) (result *gorest.OauthToken, status bool, err error)
}

// DBChannel represent channel repository (postgresql) in DB interface
type DBChannel interface {
	FindByUsername(db *gorm.DB, username string) (channel *gorestdb.OauthChannel, err error)
}

// RDSRefresh represent refresh repository (redis) in Redis interface
type RDSRefresh interface {
	Find(rdb *redis.Client, key string, index string) (string, error)
	Save(rdb *redis.Client, data interface{}, key, index string, lifetime uint) error
}

type TokenGenerator interface {
	InitToken(expired uint, ttp, username, password string) (tokenType string, expiredAccess uint)
	SetToken() (accessToken string, refreshToken string, err error)
	CheckToken(client *redis.Client, refreshPrefix string, request interface{}) (status bool, err error)
	SaveToken(client *redis.Client, accessPrefix string, refreshPrefix string) error
	DeleteToken(client *redis.Client, accessPrefix string, refreshPrefix string) error
}
