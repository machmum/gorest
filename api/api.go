package api

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/machmum/gorest/api/auth"
	al "github.com/machmum/gorest/api/auth/logging"
	at "github.com/machmum/gorest/api/auth/transport"
	"github.com/machmum/gorest/api/oauth"
	oal "github.com/machmum/gorest/api/oauth/logging"
	oat "github.com/machmum/gorest/api/oauth/transport"
	"github.com/machmum/gorest/api/product"
	prl "github.com/machmum/gorest/api/product/logging"
	prt "github.com/machmum/gorest/api/product/transport"
	"github.com/machmum/gorest/api/profile"
	pl "github.com/machmum/gorest/api/profile/logging"
	pt "github.com/machmum/gorest/api/profile/transport"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/middleware/logging"
	"github.com/machmum/gorest/middleware/secure"
	"github.com/machmum/gorest/utl/server"
	"github.com/machmum/gorest/utl/zplog"
	"net/http"
	"time"
)

// Custom errors
var (
	ErrFailedConnMysql = errors.New("Failed connecting to database")
	ErrFailedConnRedis = errors.New("Failed connecting to redis")
	ErrMethodNotFound  = errors.New("method not found")
)

// use new model
func Start(cfg *config.Configuration) {
	// Open a connection
	// store *DB to config
	db, err := gorm.Open("postgres", cfg.DB.PSN)
	if err != nil {
		panic(ErrFailedConnMysql)
	}
	defer db.Close()

	// Open redis connection
	// don't need to close connection, go-redis will do
	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Host + cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		PoolTimeout:  30 * time.Second,
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		panic(ErrFailedConnRedis)
	}

	// configure custom log for gorm
	if cfg.Debug {
		// db.SetLogger(gorm.Logger{log.New()})
		db.LogMode(cfg.Debug)
	}

	// register
	// new logger
	lgn := zplog.New()

	// Initialize echo
	// begin
	e := echo.New()
	e.Debug = cfg.Debug // read debug from config

	// register
	// middleware
	e.Use(
		middleware.Recover(),
		secure.CORS(),
		secure.Headers(),
		logging.MiddlewareLogging(lgn), // access_log
	)

	// register
	// new http error handler
	e.HTTPErrorHandler = customHTTPErrorHandler

	// register
	// request validator
	e.Validator = &server.CustomValidator{V: validator.New()}

	// register routes
	// group
	v1 := e.Group("/v1")

	// middleware validate token
	mw := secure.AuthMiddleware(rdb, cfg.Redis.Prefix.Access)

	oat.NewHTTP(oal.NewLogger(oauth.Initialize(db, rdb, cfg.Redis), lgn), v1)
	at.NewHTTP(al.NewLogger(auth.Initialize(db, rdb, cfg.Redis), lgn), v1, mw)
	pt.NewHTTP(pl.NewLogger(profile.Initialize(db, rdb, cfg.Redis), lgn), v1, mw)
	prt.NewHTTP(prl.NewLogger(product.Initialize(db, rdb, cfg.Redis), lgn), v1, mw)

	// to use tcp server
	// prepare server
	s := &http.Server{
		Addr:         cfg.Server.Port,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	}

	e.Logger.Fatal(e.StartServer(s))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	switch err.(type) {
	default:
		// default is error handle by service

	case validator.ValidationErrors:
		// error for fail validation tag

	case *echo.HTTPError:
		parseError, ok := err.(*echo.HTTPError).Internal.(*json.UnmarshalTypeError)
		if ok {
			// error for fail validation / type
			err = errors.New(parseError.Error())

		} else {
			parseError, ok := err.(*echo.HTTPError).Internal.(*json.InvalidUnmarshalError)
			if ok {
				// error for invalid unmarshal
				err = errors.New(parseError.Error())
			} else {
				// error for method / routes not found
				err = ErrMethodNotFound
				_ = server.NotFound(c, err)
				return
			}

		}
	}

	_ = server.ResponseFail(c, err)
}
