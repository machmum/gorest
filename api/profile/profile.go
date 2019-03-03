package profile

import (
	"github.com/labstack/echo"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/server"
	_ "image/gif"
	_ "image/jpeg"
	"runtime"
	"sort"
)

var (
	// width x height
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

func (p *Profile) Profile(c echo.Context, cturl string) (res gorest.ProfileSlice, err error) {
	var (
		pss     gorestdb.ProfileSimpleSlice
		profile []gorest.Profile
	)

	// leave error trace
	defer func() {
		if err != nil {
			trace(c)
		}
	}()

	pss, err = p.platform.Profile.FindByCategoryReturnSimple(p.conn.postgresql, cturl)
	if err != nil {
		return nil, err
	}

	// prolist: profile list
	prolist := make(map[uint]gorest.Profile)

	for i := range pss {

		if w, h, err = server.GetImage(pss[i].ProductImage); err != nil {
			return nil, err
		}

		prolist[pss[i].ID] = gorest.Profile{
			ID:        pss[i].ID,
			FirstName: pss[i].FirstName,
			LastName:  pss[i].LastName,
			Username:  pss[i].Username,
			Products: append(prolist[pss[i].ID].Products, gorest.Product{
				ID:   pss[i].ProductID,
				Name: pss[i].ProductName,
				Price: &gorest.Price{
					Normal: pss[i].NormalPrice,
					Sale:   pss[i].SalePrice,
				},
				Image: &gorest.ImageProduct{
					Thumbnail: &gorest.ImageRes{
						URL: pss[i].ProductImage,
						Size: &gorest.WHReq{
							Width:  w,
							Height: h,
						},
					},
				},
			}),
		}
	}

	// re-index result
	var keys []int
	for k := range prolist {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, v := range keys {
		profile = append(profile, prolist[uint(v)])
	}

	return profile, err
}

func (p *Profile) ProfileSimple(c echo.Context, profileID int) (profile gorest.Profile, err error) {
	var psimple gorestdb.ProfileSimple

	// leave error trace
	defer func() {
		if err != nil {
			trace(c)
		}
	}()

	// get profile
	if psimple, err = p.platform.Profile.FindByProfileIDReturnSimple(p.conn.postgresql, profileID); err != nil {
		return profile, err
	}

	profile = gorest.Profile{
		ID:        psimple.ID,
		FirstName: psimple.FirstName,
		LastName:  psimple.LastName,
		Username:  psimple.Username,
	}

	return profile, err
}

func (p *Profile) Product(c echo.Context, profileID int, pid int) (product gorest.Product, err error) {
	product, err = p.internal.Product.Product(c, profileID, pid)
	if err != nil {
		return product, err
	}

	return product, err
}
