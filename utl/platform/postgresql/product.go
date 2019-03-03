package postgresql

import (
	"github.com/jinzhu/gorm"
	"github.com/machmum/gorest/utl/model/postgresql"
)

var (
	ProductTable = "product"
)

type ProductDB interface {
	FindByProductID(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error)
}

// NewUser returns a new user database instance
func NewProduct() *Product {
	return &Product{}
}

// User represents the client for user table
type Product struct {
	// this field required when process the query
	field, from, join, where, sort string

	// bind is slice string to hold placeholder value
	bind []interface{}

	// gp: group_by, ob: order_by
	gp, ob string

	// flag
	// af : added_field
	af bool
}

func (p *Product) queryProductDetail(af bool) {
	p.field = "select pr.id AS product_id, pr.name AS product_name, pr.normal_price, pr.sale_price," +
	// sub-query product_image: id+name|id+name
		" (select string_agg(concat_ws('+', pi.ID, pi.name), '|')" +
		"	from product_image pi" +
		"	where pi.product_id = pr.id and pi.status = ?" +
		" ) as product_image"
	if af {
		p.field += ","
	}
	p.from = " from " + ProductTable + " pr"
	p.where = " where pr.status = ?"
	p.bind = []interface{}{StatusActive, StatusActive}
}

func (p *Product) FindByProductID(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error) {
	var q string

	p.queryProductDetail(false)
	if pid == 0 {
		p.where += " and pr.profile_id = ?"
		p.bind = append(p.bind, profileID)
		q = p.field + p.from + p.where + " order by pr.id desc limit 1"

	} else {
		p.where += " and pr.id = ?"
		p.bind = append(p.bind, pid)
		q = p.field + p.from + p.where
	}

	if err = db.Raw(q, p.bind...).Scan(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = ErrNotFoundProfile
		}
		return result, err
	}

	return result, err
}
