package oauth_test

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/machmum/gorest/api/oauth"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/mock"
	"github.com/machmum/gorest/utl/mock/mockdb"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/rds"
	"github.com/machmum/gorest/utl/secure"
	"testing"
)

// func BenchmarkTokenize(b *testing.B) {
// 	b.ResetTimer()
//
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:         "127.0.0.1:6379",
// 		Password:     "ZjY0ZWM0MDVhYjM2NmY4MTIzNGQ5OTFi",
// 		PoolTimeout:  30 * time.Second,
// 		IdleTimeout:  10 * time.Second,
// 		ReadTimeout:  30 * time.Second,
// 		WriteTimeout: 30 * time.Second,
// 	})
// 	_, err := rdb.Ping().Result()
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}
//
// 	// initiate new config
// 	cfg := &config.RDB{
// 		Lifetime: config.RDBLifetime{
// 			Token: 7,
// 			Apps:  10,
// 		},
// 		Prefix: config.RDBPrefix{
// 			Access:  "ce_access_",
// 			Refresh: "ce_refresh_",
// 			Apps:    "ce_apps",
// 		},
// 	}
//
// 	channelDB := &mockdb.Channel{
// 		FindByUsernameFn: func(db *gorm.DB, username string) (channel *gorestdb.OauthChannel, err error) {
// 			return &gorestdb.OauthChannel{
// 				ID:        1,
// 				ChannelID: 1,
// 				Version:   "1",
// 				Username:  "android",
// 				Password:  "ffbaad722ef611e6b721",
// 			}, nil
// 		},
// 	}
//
// 	for n := 0; n < b.N; n++ {
// 		s := oauth.New(nil, rdb, cfg, channelDB, rds.NewRefresh())
// 		_, _, _ = s.Tokenize(nil, "password", "android", "ffbaad722ef611e6b721", "")
// 	}
// }

func TestTokenize(t *testing.T) {
	// initiate mini-redis
	mr, err := miniredis.Run()
	if err != nil {
		logrus.Fatal("failed run miniredis")
	}
	defer mr.Close()
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	if _, err = rdb.Ping().Result(); err != nil {
		logrus.Fatal("failed redis connection")
	}

	// initiate new config
	cfg := &config.RDB{
		Lifetime: config.RDBLifetime{
			Token: 60,
			Apps:  60,
		},
		Prefix: config.RDBPrefix{
			Access:  "ce_access_",
			Refresh: "ce_refresh_",
			Apps:    "ce_apps",
		},
	}

	DBChannel := &gorestdb.OauthChannel{
		ChannelID: 4,
		Version:   "2.0",
		Username:  "android",
		Password:  "f11android",
	}

	token := secure.NewToken()
	ttps, exps := token.InitToken(20, "Bearer", "android", "f11android")
	tac, trf, err := token.SetToken()

	wantToken := &gorest.OauthToken{
		AccessToken:  tac,
		RefreshToken: trf,
		TokenType:    ttps,
		ExpireAccess: exps,
	}

	cases := []struct {
		name      string
		username  string
		password  string
		dBChannel *mockdb.Channel
		token     *mock.Token
		gtp       string
		err       error
		wantError bool
		wantToken *gorest.OauthToken
	}{

		{
			name:      "Failed to authenticate username-password",
			username:  "android",
			password:  "f11android",
			gtp:       "password",
			err:       mysql.ErrNotFoundChannel,
			wantError: true,
			dBChannel: &mockdb.Channel{
				FindByUsernameFn: func(db *gorm.DB, username string) (channel *gorestdb.OauthChannel, err error) {
					return nil, mysql.ErrNotFoundChannel
				},
			},
			token: nil,
		},
		{
			name:      "Success get token",
			username:  "android",
			password:  "f11android",
			gtp:       "password",
			err:       nil,
			wantError: false,
			dBChannel: &mockdb.Channel{
				FindByUsernameFn: func(db *gorm.DB, username string) (channel *gorestdb.OauthChannel, err error) {
					return DBChannel, nil
				},
			},
			token: &mock.Token{
				InitTokenFn: func(exp uint, ttp, u, p string) (string, uint) {
					return ttps, exps
				},
				SetTokenFn: func() (accessToken string, refreshToken string, err error) {
					return tac, trf, nil
				},
				SaveTokenFn: func(client *redis.Client, accessPrefix string, refreshPrefix string) error {
					return nil
				},
			},
			wantToken: wantToken,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := oauth.New(nil, rdb, cfg, tt.dBChannel, rds.NewRefresh(), tt.token)
			r, _, err := s.Tokenize(nil, tt.gtp, tt.username, tt.password, "")

			if tt.wantToken != nil {
				assert.Equal(t, tt.wantToken, r)
			}
			assert.Equal(t, tt.err, err)
		})
	}
}
