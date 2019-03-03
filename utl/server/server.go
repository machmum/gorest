package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/middleware"

	"github.com/labstack/echo"
)

// Used as mock server for unit test

var ErrMethodNotFound = errors.New("method not found !")

// New instantates new Echo server
func New() *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		// secure.CORS(),
		// secure.Headers(),
	)

	// register validator
	e.Validator = &CustomValidator{V: validator.New()}

	// register error
	custErr := &customErrHandler{e: e}
	e.HTTPErrorHandler = custErr.handler

	// register bind
	e.Binder = &CustomBinder{b: &echo.DefaultBinder{}}

	return e
}

// Config represents server specific config
type Config struct {
	Port                string
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	Debug               bool
}

// Start starts echo server
func Start(e *echo.Echo, cfg *Config) {
	s := &http.Server{
		Addr:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
	}
	e.Debug = cfg.Debug

	// Start server
	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

type customErrHandler struct {
	e *echo.Echo
}

func (ce *customErrHandler) handler(err error, c echo.Context) {
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
				_ = NotFound(c, err)
				return
			}

		}
	}

	_ = ResponseFail(c, err)
}
