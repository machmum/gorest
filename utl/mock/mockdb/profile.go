package mockdb

import (
	"github.com/jinzhu/gorm"
	"github.com/machmum/gorest/utl/model/postgresql"
)

// Profile database mock
type Profile struct {
	FindByUsernameFn              func(db *gorm.DB, username string) (result gorestdb.Profile, err error)
	FindByProfileIDReturnSimpleFn func(db *gorm.DB, id int) (result gorestdb.ProfileSimple, err error)
	FindByCategoryReturnSimpleFn  func(db *gorm.DB, category string) (slice gorestdb.ProfileSimpleSlice, err error)
}

// FindByUsername mock
func (p *Profile) FindByUsername(db *gorm.DB, username string) (gorestdb.Profile, error) {
	return p.FindByUsernameFn(db, username)
}

func (p *Profile) FindByProfileIDReturnSimple(db *gorm.DB, id int) (result gorestdb.ProfileSimple, err error) {
	return p.FindByProfileIDReturnSimpleFn(db, id)
}

func (p *Profile) FindByCategoryReturnSimple(db *gorm.DB, cturl string) (result gorestdb.ProfileSimpleSlice, err error) {
	return p.FindByCategoryReturnSimpleFn(db, cturl)
}
