package transport

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/profile"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/server"
	"strings"
)

var (
	user = "user"
)

// HTTP represents auth http service
type HTTP struct {
	svc profile.Service
}

func NewHTTP(svc profile.Service, er *echo.Group, mw echo.MiddlewareFunc) {
	h := HTTP{svc: svc}

	er.POST("/profile", h.profile, mw)
	er.POST("/profile/detail", h.profileDetail, mw)
}

func (h *HTTP) profile(c echo.Context) error {
	var (
		err    error
		req    = new(gorest.ProfileReq)
		r      interface{} // hold result
		result = make(map[string]interface{}, 2)
	)

	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	// check scope
	scope := server.NewScope(strings.Split(req.Scope, "+"), server.ScopeProfile)
	if err = scope.Check(); err != nil {
		return err
	}

	if scope.ScopeSet["profile"] == true {
		// do user
		if req.Category == user {
			r, err = h.svc.Profile(c, req.Category)
			if err != nil {
				return err
			}
		}

		result["profile"] = r
	}

	return server.ResponseOK(c, "", result)
}

func (h *HTTP) profileDetail(c echo.Context) error {
	var (
		err    error
		req    = new(gorest.PflDetailReq)
		result = make(map[string]interface{}, 3)
	)

	if err = c.Bind(req); err != nil {
		return err
	}

	if err = c.Validate(req); err != nil {
		return err
	}

	// check scope
	scope := server.NewScope(strings.Split(req.Scope, "+"), server.ScopeProfileDetail)
	if err = scope.Check(); err != nil {
		return err
	}

	// validate scope: size
	if req.Size.Profile == nil {
		return server.ErrSizeNotFound
	}

	if scope.ScopeSet["profile"] == true && req.Size.Profile == nil {
		return server.ErrSizeProfileNotFound
	}
	if scope.ScopeSet["menu"] == true && req.Size.Menu == nil {
		return server.ErrSizeMenuNotFound
	}
	if scope.ScopeSet["product"] == true && req.Size.Product == nil {
		return server.ErrSizeProductNotFound
	}

	if scope.ScopeSet["profile"] == true {
		prof, err := h.svc.ProfileSimple(c, req.ProfileID)
		if err != nil {
			return err
		}

		result["profile"] = prof
	}

	if scope.ScopeSet["product"] == true {
		product, err := h.svc.Product(c, req.ProfileID, req.ProductID)
		if err != nil {
			return err
		}

		result["product"] = product
	}

	return server.ResponseOK(c, "", result)
}
