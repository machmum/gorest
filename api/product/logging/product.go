package product

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/api/product"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/zplog"
)

var tracer string

const name = "product"

// LogService represents auth logging service
type LogService struct {
	svc    product.Service
	logger zplog.Logger
}

func NewLogger(svc product.Service, logger zplog.Logger) *LogService {
	return &LogService{
		svc:    svc,
		logger: logger,
	}
}

func (ls *LogService) Product(c echo.Context, profileID int, pid int) (products gorest.Product, err error) {
	return ls.svc.Product(c, profileID, pid)
}
