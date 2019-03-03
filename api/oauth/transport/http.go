package transport

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/oauth"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/server"
)

// HTTP represents auth http service
type HTTP struct {
	svc oauth.Service
}

func NewHTTP(svc oauth.Service, er *echo.Group) {
	h := HTTP{svc: svc}

	g := er.Group("/oauth/token")

	g.POST("/request", h.request)
	g.POST("/refresh", h.refresh)
}

func (h *HTTP) request(c echo.Context) (err error) {
	var (
		status bool
		result = new(gorest.OauthToken)
	)

	cred := new(gorest.Credentials)

	err = c.Bind(cred)
	if err != nil {
		return err
	}

	err = c.Validate(cred)
	if err != nil {
		return err
	}

	gtype := "request"

	if result, status, err = h.svc.Tokenize(c, gtype, cred.Username, cred.Password, ""); err != nil {
		if status {
			return server.ResponseUnauthorized(c, err)
		}

		return err
	}

	return server.ResponseOK(c, "", result)
}

func (h *HTTP) refresh(c echo.Context) (err error) {
	var (
		status bool
		result = new(gorest.OauthToken)
	)

	cred := new(gorest.CredentialsRefresh)

	err = c.Bind(cred)
	if err != nil {
		return err
	}

	err = c.Validate(cred)
	if err != nil {
		return err
	}

	gtype := "refresh"

	if result, status, err = h.svc.Tokenize(c, gtype, cred.Username, cred.Password, cred.RefreshToken); err != nil {
		if status {
			return server.ResponseUnauthorized(c, err)
		}

		return err
	}

	return server.ResponseOK(c, "", result)
}
