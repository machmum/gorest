package secure

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/gofrs/uuid"
	"strings"
	"time"
)

var (
	ErrInvalidToken     = errors.New("token invalid or expired")
	ErrNotFoundToken    = errors.New("token not found")
	ErrNotFoundReqToken = errors.New("refresh token not found")
	ErrInvalidType      = errors.New("invalid type, type must be a map")
	ErrCekToken         = errors.New("invalid parameter used to check token")
	ErrNilValue         = errors.New("map value is nil")

	tokenLen = 6
)

type TokenSource interface {
	InitToken(expired uint, ttp, username, password string) (tokenType string, expiredAccess uint)
	SetToken() (accessToken string, refreshToken string, err error)
	CheckToken(client *redis.Client, refreshPrefix string, request interface{}) (status bool, err error)
	SaveToken(client *redis.Client, accessPrefix string, refreshPrefix string) error
	DeleteToken(client *redis.Client, accessPrefix string, refreshPrefix string) error
}

type (
	// New object token
	Token struct {
		// AccessToken is the token that authorizes and authenticates
		// the requests.
		AccessToken string `json:"access_token"`

		// TokenType is the type of token.
		// The Type method returns either this or "Bearer", the default.
		TokenType string `json:"token_type,omitempty"`

		// RefreshToken is a token that's used by the application
		// (as opposed to the user) to refresh the access token
		// if it expires.
		RefreshToken string `json:"refresh_token,omitempty"`

		// Expiry is the optional expiration time of the access token.
		//
		// If zero, TokenSource implementations will reuse the same
		// token forever and RefreshToken or equivalent
		// mechanisms for that TokenSource will not be used.
		ExpiryAccess time.Duration `json:"expiry_access,omitempty"`

		// Expiry is the optional expiration time of the access token.
		//
		// If zero, TokenSource implementations will reuse the same
		// token forever and RefreshToken or equivalent
		// mechanisms for that TokenSource will not be used.
		ExpiryRefresh time.Duration `json:"expiry_refresh,omitempty"`

		// Saved token
		// SavedAccess is access token saved in redis
		// SavedRefresh is refresh token saved in redis
		// used for authenticating refresh token
		SavedAccess  string `json:"saved_access,omitempty"`
		SavedRefresh string `json:"saved_refresh,omitempty"`

		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
)

func NewToken() *Token {
	return &Token{}
}

func (t *Token) InitToken(expired uint, ttp, username, password string) (tokenType string, expiredAccess uint) {
	t.Username = username
	t.Password = password
	t.TokenType = ttp
	t.ExpiryAccess = time.Duration(expired) * time.Second

	return t.TokenType, expired
}

func (t *Token) SetToken() (at, rt string, err error) {
	var (
		str string
	)

	// set token type
	// t.TokenType = tokenType

	// begin create token
	// token is uuid v4 : its return randomly generated UUID.
	for i := 0; i < 2; i++ {
		err = func(count int, s string) error {
			// create UUID based on tokenLen length
			for j := 0; j < tokenLen; j++ {
				// initiate new UUID
				u, err := uuid.NewV4()
				if err != nil {
					return err
				}
				s += uuid.Must(u, err).String() + "-"
			}
			if count == 0 {
				t.AccessToken = strings.TrimSuffix(s, "-")
			} else {
				t.RefreshToken = strings.TrimSuffix(s, "-")
			}
			return nil
		}(i, str)

		if err != nil {
			break
		}
	}

	return t.AccessToken, t.RefreshToken, err
}

// Check token check given refresh token with refresh token in redis
//
// status : false, undefined error found
// status : true, defined error found
func (t *Token) CheckToken(conn *redis.Client, rpr string, request interface{}) (status bool, err error) {
	var (
		username string
		password string
	)

	// check request
	// return error if request is not map[string]interface{}
	val, ok := request.(map[string]interface{})
	if !ok {
		return false, ErrInvalidType
	}
	if val == nil {
		return false, ErrNilValue
	}

	// assign request body
	// request body contains username, password and refresh_token
	for key, val := range request.(map[string]interface{}) {

		if key == "username" {
			username = val.(string)
		} else if key == "password" {
			password = val.(string)
		} else if key == "refresh_token" {
			t.SavedRefresh = val.(string)
		}
	}

	// return error if refresh_token not set in request
	if t.SavedRefresh == "" {
		return true, ErrNotFoundReqToken
	}

	// get token saved in redis
	// get by given refresh_token
	trds, err := conn.HGetAll(rpr + t.SavedRefresh).Result()
	if err != nil {
		return false, err
	}

	if len(trds) == 0 {
		return true, ErrInvalidToken

	} else {

		// json to struct
		if err = json.Unmarshal([]byte(trds["token"]), &t); err != nil {
			return false, err
		}

		// validate given username with saved username in refresh_token
		// validate given password with saved password in refresh_token
		if username == t.Username && password == t.Password {

			// save saved access_token
			t.SavedAccess = t.AccessToken

		} else {
			return true, ErrCekToken
		}
	}

	return false, nil
}

func (t *Token) SaveToken(conn *redis.Client, apr, rpr string) (err error) {
	prefixes := map[string]string{
		"access":  apr + t.AccessToken,
		"refresh": rpr + t.RefreshToken,
	}

	// set expiry refresh 2 times expiry access
	t.ExpiryRefresh = t.ExpiryAccess * 2

	// populate refresh object to redis
	refresh, _ := json.Marshal(map[string]interface{}{
		"access_token":   t.AccessToken,
		"token_type":     t.TokenType,
		"expiry_refresh": t.ExpiryRefresh / time.Second,
		"username":       t.Username,
		"password":       t.Password,
	})

	// populate access object to redis
	access, _ := json.Marshal(map[string]interface{}{
		"refresh_token": t.RefreshToken,
		"token_type":    t.TokenType,
		"expiry_access": t.ExpiryAccess / time.Second,
		"username":      t.Username,
	})

	for key, prefix := range prefixes {
		var p string

		err, p = func(prefix string, key string) (error, string) {
			if key == "access" {
				if _, err = conn.HMSet(prefix, map[string]interface{}{"token": access}).Result(); err != nil {
					return err, prefix
				}
				if _, err := conn.Expire(prefix, t.ExpiryAccess).Result(); err != nil {
					return err, prefix
				}
			} else {
				if _, err = conn.HMSet(prefix, map[string]interface{}{"token": refresh}).Result(); err != nil {
					return err, prefix
				}
				if _, err := conn.Expire(prefix, t.ExpiryRefresh).Result(); err != nil {
					return err, prefix
				}
			}
			return nil, ""
		}(prefix, key)

		if err != nil {
			_, _ = conn.HDel(p, "token").Result() // force delete previous token
			break
		}
	}

	return err
}

func (t *Token) DeleteToken(conn *redis.Client, apr, rpr string) (err error) {
	var (
		oldkey string
		newkey string
	)

	// validate
	if t.SavedAccess == "" || t.SavedRefresh == "" {
		return ErrNotFoundToken
	}

	// delete old access token
	oldkey = apr + t.SavedAccess
	if _, err = conn.HDel(oldkey, "token").Result(); err != nil {
		return err
	}

	// move refresh token
	oldkey = rpr + t.SavedRefresh
	newkey = rpr + t.RefreshToken
	if _, err = conn.Rename(oldkey, newkey).Result(); err != nil {
		return err
	}

	return err
}
