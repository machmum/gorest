package secure_test

import (
	"errors"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/secure"
	"testing"
	"time"
)

func TestSetToken(t *testing.T) {

	token := secure.NewToken()
	ttps, exps := token.InitToken(20, "Bearer", "android", "f11android")
	tac, trf, err := token.SetToken()
	if err != nil {
		panic(err)
	}

	wantToken := &gorest.OauthToken{
		AccessToken:  tac,
		RefreshToken: trf,
		TokenType:    ttps,
		ExpireAccess: exps,
	}

	assert.NotEqual(t, token, wantToken)
}

func TestCheckToken(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		logrus.Fatal("failed run miniredis")
	}
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{Addr: s.Addr()})
	if _, err = rdb.Ping().Result(); err != nil {
		logrus.Fatal("failed redis connection")
	}

	token := secure.NewToken()
	_, _ = token.InitToken(20, "Bearer", "android", "f11android")
	tac, trf, err := token.SetToken()
	if err != nil {
		panic(err)
	}

	if err = token.SaveToken(rdb, "ce_access_", "ce_refresh_"); err != nil {
		logrus.Fatal("failed save token")
	}

	cases := []struct {
		name    string
		request map[string]interface{}
		token   *secure.Token
		errname error
	}{
		{
			name: "Failed username differ",
			request: map[string]interface{}{
				"username":      "web",
				"password":      "ffbaad722ef611e6b721",
				"refresh_token": trf,
			},
			token: &secure.Token{
				AccessToken:  tac,
				RefreshToken: trf,
				TokenType:    "Bearer",
				ExpiryAccess: time.Duration(60 * time.Second),
				Username:     "android",
				Password:     "f11android",
			},
			errname: errors.New("invalid parameter used to check token"),
		},
		{
			name: "Failed refresh token not found",
			request: map[string]interface{}{
				"username":      "android",
				"password":      "ffbaad722ef611e6b721",
				"refresh_token": "",
			},
			token: &secure.Token{
				AccessToken:  "",
				RefreshToken: trf,
				TokenType:    "Bearer",
				ExpiryAccess: time.Duration(60 * time.Second),
				Username:     "android",
				Password:     "f11android",
			},
			errname: errors.New("refresh token not found"),
		},
		{
			name: "Success check token",
			request: map[string]interface{}{
				"username":      "android",
				"password":      "f11android",
				"refresh_token": trf,
			},
			token: &secure.Token{
				AccessToken:  tac,
				RefreshToken: trf,
				TokenType:    "Bearer",
				ExpiryAccess: time.Duration(60 * time.Second),
				Username:     "android",
				Password:     "f11android",
			},
			errname: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, err = token.CheckToken(rdb, "ce_refresh_", tt.request)

			assert.Equal(t, tt.errname, err)
		})
	}
}
