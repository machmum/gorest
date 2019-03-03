package auth

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/auth"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/zplog"
	"time"
)

var tracer string

const name = "auth"

// LogService represents auth logging service
type LogService struct {
	svc    auth.Service
	logger zplog.Logger
}

func NewLogger(svc auth.Service, logger zplog.Logger) *LogService {
	return &LogService{
		svc:    svc,
		logger: logger,
	}
}

func (ls *LogService) Login(c echo.Context, user string, pass string) (login gorest.Profile, err error) {
	defer func(begin time.Time) {
		e := c.Get("etrace")
		if e != nil {
			tracer = zplog.Trace(e)
		}

		ls.logger.Log(name, "Login request", err,
			map[string]interface{}{
				"req":   map[string]string{"user": user, "pass": pass},
				"res":   login,
				"took":  time.Since(begin),
				"trace": tracer,
			},
		)
	}(time.Now())
	return ls.svc.Login(c, user, pass)
}

func (ls *LogService) Logout(c echo.Context, token string) (err error) {
	defer func(begin time.Time) {
		e := c.Get("etrace")
		if e != nil {
			tracer = zplog.Trace(e)
		}

		ls.logger.Log(name, "Logout request", err,
			map[string]interface{}{
				"took":  time.Since(begin),
				"trace": tracer,
			},
		)
	}(time.Now())
	return ls.svc.Logout(c, token)
}
