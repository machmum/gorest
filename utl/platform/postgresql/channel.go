package postgresql

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/machmum/gorest/utl/model/postgresql"
)

var (
	ChannelTable       = "oauth_channel"
	ErrNotFoundChannel = errors.New("failed to get channel")

	StatusActive = 1
)

// NewUser returns a new user database instance
func NewChannel() *Channel {
	return &Channel{}
}

// User represents the client for user table
type Channel struct{}

func (c *Channel) FindByUsername(db *gorm.DB, username string) (*gorestdb.OauthChannel, error) {
	var result = new(gorestdb.OauthChannel)

	q := "select channel_id, username, password from " + ChannelTable + " where username = ?"

	if err := db.Raw(q, username).Scan(&result).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = ErrNotFoundChannel
		}
		return nil, err
	}

	return result, nil
}
