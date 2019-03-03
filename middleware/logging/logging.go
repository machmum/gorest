package logging

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/utl/zplog"
	"time"
)

func MiddlewareLogging(logger zplog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Log("", "Incoming request", nil,
				map[string]interface{}{
					"at":         time.Now().Format("2006-01-02 15:04:05"),
					"method":     c.Request().Method,
					"uri":        c.Request().URL.String(),
					"ip":         c.Request().RemoteAddr,
					"host":       c.Request().Host,
					"user_agent": c.Request().UserAgent(),
					"code":       c.Response().Status,
				},
			)

			return next(c)
		}
	}
}
