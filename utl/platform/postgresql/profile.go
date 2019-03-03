package postgresql

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/machmum/gorest/utl/model/postgresql"
)

var (
	ProfileTable      = "profile"
	ProductImageTable = "product_image"

	ErrNotFoundProfile = errors.New("profile not found")
)

// NewUser returns a new user database instance
func NewProfile() *Profile {
	return &Profile{}
}

type ProfileDB interface {
	FindByUsername(db *gorm.DB, username string) (result gorestdb.Profile, err error)
	FindByProfileIDReturnSimple(db *gorm.DB, id int) (result gorestdb.ProfileSimple, err error)
	FindByCategoryReturnSimple(db *gorm.DB, category string) (slice gorestdb.ProfileSimpleSlice, err error)
}

// User represents the client for user table
type Profile struct {
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

// InitProfileQuery used to start query profile
// by default, get profile active / p.status = 1
//
// returned field:
//
// profile.*
// referential_profile.name as status_profile 		where 		key = category. 	table_name = profile. field_name = category
// referential_profile.name as status_verification 	where 		key = verification. table_name = profile. field_name =  verification
// referential_profile.name as status_popular 		where 		key = popular. 		table_name = profile. field_name = popular
//
// condition: 										where 		profile.status = 1
func (p *Profile) profileFull(af bool) {
	p.field = "select p.* "
	if af {
		p.field += ", "
	}
	p.from = " from " + ProfileTable + " p"
	p.join = " left join "
	p.where = " where p.`status` = ?"
	p.bind = []interface{}{StatusActive}
}

// Returned:
//
// profile: id. first_name. last_name. username. image
//
// condition: where profile.status = 1
func (p *Profile) profileSimple(af bool) {
	p.field = "select p.id, p.first_name, p.last_name, p.username"
	if af {
		p.field += ","
	}
	p.from = " from " + ProfileTable + " p"
	p.join = " left join "
	p.where = " where p.`status` = ?"
	p.bind = []interface{}{StatusActive}
}

// Find profile by username
// returned: profileFull
//
// conditions: 	where username = ?
// limit 1
func (p *Profile) FindByUsername(db *gorm.DB, username string) (result gorestdb.Profile, err error) {
	p.profileFull(false)

	q := p.field + p.from + " where username = ? limit 1"
	bind := []interface{}{username}

	if err = db.Raw(q, bind...).Scan(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = ErrNotFoundProfile
		}
		return result, err
	}

	return result, err
}

// Find profile by profile_id
// returned: profileSimple
//
func (p *Profile) FindByProfileIDReturnSimple(db *gorm.DB, id int) (result gorestdb.ProfileSimple, err error) {
	p.profileSimple(false)

	q := p.field + p.from + " where p.id = ? limit 1"
	bind := []interface{}{id}

	if err = db.Raw(q, bind...).Scan(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = ErrNotFoundProfile
		}
		return result, err
	}

	return result, err
}

// Find Profile by category
// returned field:
//
// profileSimple
// product: 		id as product_id. name as product_name. normal_price. sale_price
// product_image: 	id as pi_id. name as product_image
//
// condition:
// where referential_profile.key
func (p *Profile) FindByCategoryReturnSimple(db *gorm.DB, category string) (slice gorestdb.ProfileSimpleSlice, err error) {
	p.profileSimple(true)

	q := p.field +
		" pr.id as product_id, pr.name as product_name, pr.normal_price, pr.sale_price, pi.id as pi_id, pi.name as product_image" +
		p.from +
		" left join " + ProductTable + " pr on pr.profile_id = p.id" +
		" left join " + ProductImageTable + " pi on pi.product_id = pr.id" +
		" where pr.status = ? and pi.status = ?"
	bind := []interface{}{StatusActive, StatusActive}

	if err = db.Raw(q, bind...).Scan(&slice).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = ErrNotFoundProfile
		}
		return nil, err
	}

	return slice, nil
}
