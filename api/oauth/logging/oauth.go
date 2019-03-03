package oauth

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/oauth"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/zplog"
	"time"
)

var tracer string

const name = "oauth"

// LogService represents auth logging service
type LogService struct {
	svc    oauth.Service
	logger zplog.Logger
}

func NewLogger(svc oauth.Service, logger zplog.Logger) *LogService {
	return &LogService{
		svc:    svc,
		logger: logger,
	}
}

func (ls *LogService) Tokenize(c echo.Context, gtype, u, p, rt string) (oauthToken *gorest.OauthToken, status bool, err error) {
	defer func(begin time.Time) {
		e := c.Get("etrace")
		if e != nil {
			tracer = zplog.Trace(e)
		}

		ls.logger.Log(name, "Login request", err,
			map[string]interface{}{
				"req":   map[string]string{"grand_type": gtype, "user": u, "pass": p, "refresh_token": rt},
				"took":  time.Since(begin),
				"trace": tracer,
			},
		)
	}(time.Now())
	return ls.svc.Tokenize(c, gtype, u, p, rt)
}
