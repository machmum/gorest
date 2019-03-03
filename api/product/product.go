package product

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/server"
	_ "image/gif"
	_ "image/jpeg"
	"runtime"
	"strings"
)

var (
	// redis attributes
	w, h int
)

func trace(c echo.Context) {
	if c != nil {
		function, file, line, _ := runtime.Caller(1)
		etrace := map[string]interface{}{
			"func": function,
			"file": file,
			"line": line,
		}
		c.Set("etrace", etrace)
	}
}

func (p *Product) Product(c echo.Context, profileID int, pid int) (product gorest.Product, err error) {
	var (
		j int
	)

	// leave error trace
	defer func() {
		if err != nil {
			trace(c)
		}
	}()

	prdb, err := p.platform.Product.FindByProductID(p.conn.postgresql, profileID, pid)
	if err != nil {
		return product, err
	}

	tmpImg := make(gorest.ImageProductSlice, 2, 2)

	pi := strings.Split(prdb.ProductImage, "|")
	for i := range pi {
		source := strings.Split(pi[i], "+")[1]
		if j == 0 {
			w, h, err = server.GetImage(source)
			if err != nil {
				break
			}
		}
		j++

		tmpImg[i] = gorest.ImageProduct{
			Thumbnail: &gorest.ImageRes{
				URL: source,
				Size: &gorest.WHReq{
					Width:  w,
					Height: h,
				},
			},
		}
	}

	product = gorest.Product{
		ID:   prdb.ProductID,
		Name: prdb.ProductName,
		Price: &gorest.Price{
			Normal: prdb.NormalPrice,
			Sale:   prdb.SalePrice,
		},
		Images: tmpImg,
	}

	return product, err
}
