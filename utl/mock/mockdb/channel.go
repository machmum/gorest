package mockdb

import (
	"github.com/jinzhu/gorm"
	"github.com/machmum/gorest/utl/model/postgresql"
)

// User database mock
type Channel struct {
	FindByUsernameFn func(db *gorm.DB, username string) (channel *gorestdb.OauthChannel, err error)
}

// FindByUsername mock
func (c *Channel) FindByUsername(db *gorm.DB, username string) (*gorestdb.OauthChannel, error) {
	return c.FindByUsernameFn(db, username)
}
