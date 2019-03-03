package auth_test

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/machmum/gorest/api/auth"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/mock/mockdb"
	"github.com/machmum/gorest/utl/mock/mockrds"
	"github.com/machmum/gorest/utl/platform/rds"
	"testing"
)

func TestLogin(t *testing.T) {
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

	cases := []struct {
		name      string
		username  string
		password  string
		profileDB *mockdb.Profile
		gtp       string
		err       error
		wantError bool
	}{
		{
			name:     "Wrong password",
			username: "onedeca",
			password: "ffbaad722ef611e6b721",
			// err:       mysql.ErrNotFoundProfile,
			wantError: true,
			profileDB: &mockdb.Profile{
				// FindByUsernameFn: func(db *gorm.DB, username string) (profile *gorestdb.Profile, err error) {
				// 	return nil, mysql.ErrNotFoundProfile
				// },
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, rdb, cfg, tt.profileDB, rds.NewRefresh())
			_, err := s.Login(nil, tt.username, tt.password)

			assert.Equal(t, tt.err, err)

		})
	}
}

func TestLogout(t *testing.T) {
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

	cases := []struct {
		name       string
		profileDB  *mockdb.Profile
		refreshRDS *mockrds.Refresh
		err        error
	}{
		{
			name: "Success logout",
			err:  nil,
			refreshRDS: &mockrds.Refresh{
				DeleteFn: func(rdb *redis.Client, key, index string) error {
					return nil
				},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, rdb, cfg, tt.profileDB, tt.refreshRDS)
			err := s.Logout(nil, "test")

			assert.Equal(t, tt.err, err)

		})
	}
}
