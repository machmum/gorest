package transport

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/auth"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/server"
	"strings"
)

// HTTP represents auth http service
type HTTP struct {
	svc auth.Service
}

func NewHTTP(svc auth.Service, er *echo.Group, mw echo.MiddlewareFunc) {
	h := HTTP{svc: svc}

	g := er.Group("/auth")

	g.POST("/login", h.login, mw)
	g.GET("/logout", h.logout, mw)
}

func (h *HTTP) login(c echo.Context) (err error) {
	var (
		cred   = new(gorest.LoginReq)
		result = make(map[string]interface{}, 1)
	)

	if err = c.Bind(cred); err != nil {
		return err
	}

	if err = c.Validate(cred); err != nil {
		return err
	}

	// check scope
	scope := server.NewScope(strings.Split(cred.Scope, "+"), server.ScopeLogin)
	if err = scope.Check(); err != nil {
		return err
	}

	if scope.ScopeSet["profile"] == true {

		result["profile"], err = h.svc.Login(c, cred.Username, cred.Password)
		if err != nil {
			return err
		}
	}

	return server.ResponseOK(c, "", result)
}

func (h *HTTP) logout(c echo.Context) (err error) {
	err = h.svc.Logout(c, c.Get("oauth_refresh").(string))
	if err != nil {
		return err
	}

	return server.ResponseOK(c, "Success logout", nil)
}
