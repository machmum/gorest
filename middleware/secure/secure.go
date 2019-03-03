package secure

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/machmum/gorest/utl/secure"
	"github.com/machmum/gorest/utl/server"
	"strings"
)

// Headers adds general security headers for basic security measures
func Headers() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Protects from MimeType Sniffing
			c.Response().Header().Set("X-Content-Type-Options", "nosniff")
			// Prevents browser from prefetching DNS
			c.Response().Header().Set("X-DNS-Prefetch-Control", "off")
			// Denies website content to be served in an iframe
			c.Response().Header().Set("X-Frame-Options", "DENY")
			c.Response().Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
			// Prevents Internet Explorer from executing downloads in site's context
			c.Response().Header().Set("X-Download-Options", "noopen")
			// Minimal XSS protection
			c.Response().Header().Set("X-XSS-Protection", "1; mode=block")
			return next(c)
		}
	}
}

// CORS adds Cross-Origin Resource Sharing support
func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		MaxAge:           86400,
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "PATCH", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}

func AuthMiddleware(conn *redis.Client, prefix string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// return next(c)

			var token map[string]interface{}
			var rk = "oauth_refresh"

			req := c.Request()

			// get header authorization
			// add to prefix to get data redis
			auth := req.Header.Get("Authorization")
			if len(auth) < 7 {
				auth = " "
			}
			auth = strings.Replace(auth, "Bearer ", "", 1)

			key := prefix + auth

			// get token saved in redis
			// get by given refresh_token
			// return if token not found in redis
			trds, err := conn.HGetAll(key).Result()
			if err != nil {
				c.Error(err)
				return nil
			}
			if len(trds) < 1 {
				return server.ResponseUnauthorized(c, secure.ErrInvalidToken)
			}

			// json to struct
			// set refresh token
			if err := json.Unmarshal([]byte(trds["token"]), &token); err != nil {
				c.Error(err)
				return err
			}

			c.Set(rk, token["refresh_token"])

			return next(c)
		}
	}
}
