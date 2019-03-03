package gorestdb

type (
	OauthChannel struct {
		ID        uint   `json:"id,omitempty" gorm:"primary_key"`
		ChannelID int    `json:"channel_id,omitempty" gorm:"column:channel_id"`
		Version   string `json:"version,omitempty" gorm:"omitempty"`
		Username  string `json:"username,omitempty"`
		Password  string `json:"password,omitempty"`
		Base
	}
)
