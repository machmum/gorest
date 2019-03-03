package transport

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/product"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/server"
	"strings"
)

var (
	cooker = "cook"
)

// HTTP represents auth http service
type HTTP struct {
	svc product.Service
}

func NewHTTP(svc product.Service, er *echo.Group, mw echo.MiddlewareFunc) {
	h := HTTP{svc: svc}

	er.POST("/product", h.product, mw)
	// er.POST("/product/detail", h.profileDetail)
}

func (h *HTTP) product(c echo.Context) error {
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

	if scope.ScopeSet["product"] == true && req.Size.Product == nil {
		return server.ErrSizeProductNotFound
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
