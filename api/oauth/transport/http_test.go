package transport_test

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/machmum/gorest/api/oauth"
	oal "github.com/machmum/gorest/api/oauth/logging"
	"github.com/machmum/gorest/api/oauth/transport"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/mock"
	"github.com/machmum/gorest/utl/mock/mockdb"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/mysql"
	"github.com/machmum/gorest/utl/platform/rds"
	"github.com/machmum/gorest/utl/secure"
	"github.com/machmum/gorest/utl/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

	_ = &gorestdb.OauthChannel{
		// DBChannel := &gorestdb.OauthChannel{
		ChannelID: 4,
		Version:   "2.0",
		Username:  "android",
		Password:  "f11android",
	}

	token := secure.NewToken()
	ttps, exps := token.InitToken(20, "Bearer", "android", "f11android")
	tac, trf, err := token.SetToken()

	_ = &gorest.OauthToken{
		// wantToken := &gorest.OauthToken{
		AccessToken:  tac,
		RefreshToken: trf,
		TokenType:    ttps,
		ExpireAccess: exps,
	}

	cases := []struct {
		name       string
		path       string
		req        string
		dBChannel  *mockdb.Channel
		token      *mock.Token
		wantStatus int
		wantResp   *gorest.OauthToken
	}{
		{
			name:       "Invalid request",
			path:       "/v1/token",
			req:        `{"username":"juzernejm"}`,
			wantStatus: http.StatusBadRequest,
		},
		// {
		// 	name:       "Fail on FindByUsername",
		// 	req:        `{"username":"juzernejm","password":"hunter123"}`,
		// 	wantStatus: http.StatusInternalServerError,
		// 	udb: &mockdb.User{
		// 		FindByUsernameFn: func(orm.DB, string) (*gorsk.User, error) {
		// 			return nil, gorsk.ErrGeneric
		// 		},
		// 	},
		// },
		// {
		// 	name:       "Success",
		// 	req:        `{"username":"juzernejm","password":"hunter123"}`,
		// 	wantStatus: http.StatusOK,
		// 	dBChannel: &mockdb.Channel{
		// 				FindByUsernameFn: func(db *gorm.DB, username string) (channel *gorestdb.OauthChannel, err error) {
		// 					return DBChannel, nil
		// 				},
		// 			},
		// 	jwt: &mock.JWT{
		// 		GenerateTokenFn: func(*gorsk.User) (string, string, error) {
		// 			return "jwttokenstring", mock.TestTime(2018).Format(time.RFC3339), nil
		// 		},
		// 	},
		// 	sec: &mock.Secure{
		// 		HashMatchesPasswordFn: func(string, string) bool {
		// 			return true
		// 		},
		// 		TokenFn: func(string) string {
		// 			return "refreshtoken"
		// 		},
		// 	},
		// 	wantResp: &gorsk.AuthToken{Token: "jwttokenstring", Expires: mock.TestTime(2018).Format(time.RFC3339), RefreshToken: "refreshtoken"},
		// },
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := server.New()
			v1 := r.Group("/v1")

			transport.NewHTTP(oal.NewLogger(oauth.New(nil, rdb, cfg, tt.dBChannel, rds.NewRefresh(), tt.token), nil), v1)

			ts := httptest.NewServer(r)
			defer ts.Close()

			path := ts.URL + tt.path
			res, err := http.Post(path, "application/json", bytes.NewBufferString(tt.req))
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			logrus.Print(path)
			logrus.Fatal(res.StatusCode)

			response := new(gorest.OauthToken)
			if err := json.NewDecoder(res.Body).Decode(response); err != nil {
				t.Fatal(err)
			}
			m, _ := json.Marshal(response)
			t.Fatal(string(m))
			t.Fatal(response)

			if tt.wantResp != nil {
				response := new(gorest.OauthToken)
				if err := json.NewDecoder(res.Body).Decode(response); err != nil {
					t.Fatal(err)
				}
				t.Fatal(response)
				tt.wantResp.RefreshToken = response.RefreshToken
				assert.Equal(t, tt.wantResp, response)
			}
			assert.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
