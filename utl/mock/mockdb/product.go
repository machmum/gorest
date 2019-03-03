package mockdb

import (
	"github.com/jinzhu/gorm"
	"github.com/machmum/gorest/utl/model/postgresql"
)

type Product struct {
	FindByProductIDFn func(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error)
}

func (pr *Product) FindByProductID(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error) {
	return pr.FindByProductIDFn(db, profileID, pid)
}
