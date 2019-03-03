package profile

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/profile"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/zplog"
	"time"
)

var tracer string

const name = "profile"

// LogService represents auth logging service
type LogService struct {
	svc    profile.Service
	logger zplog.Logger
}

func NewLogger(svc profile.Service, logger zplog.Logger) *LogService {
	return &LogService{
		svc:    svc,
		logger: logger,
	}
}

func (ls *LogService) Profile(c echo.Context, cturl string) (result gorest.ProfileSlice, err error) {
	defer func(begin time.Time) {
		e := c.Get("etrace")
		if e != nil {
			tracer = zplog.Trace(e)
		}

		ls.logger.Log(name, "Login request", err,
			map[string]interface{}{
				"req":   map[string]interface{}{},
				"res":   result,
				"took":  time.Since(begin),
				"trace": tracer,
			},
		)
	}(time.Now())
	return ls.svc.Profile(c, cturl)
}

func (ls *LogService) ProfileSimple(c echo.Context, profileID int) (profile gorest.Profile, err error) {
	return ls.svc.ProfileSimple(c, profileID)
}

func (ls *LogService) Product(c echo.Context, profileID int, pid int) (products gorest.Product, err error) {
	return ls.svc.Product(c, profileID, pid)
}
