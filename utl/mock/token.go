package mock

import (
	"github.com/go-redis/redis"
)

type Token struct {
	InitTokenFn   func(expired uint, ttp, username, password string) (tokenType string, expiredAccess uint)
	SetTokenFn    func() (accessToken string, refreshToken string, err error)
	CheckTokenFn  func(client *redis.Client, refreshPrefix string, request interface{}) (status bool, err error)
	SaveTokenFn   func(client *redis.Client, accessPrefix string, refreshPrefix string) error
	DeleteTokenFn func(client *redis.Client, accessPrefix string, refreshPrefix string) error
}

// mock secure/token.go
func (t *Token) InitToken(expired uint, ttp, username, password string) (tokenType string, expiredAccess uint) {
	return t.InitTokenFn(expired, ttp, username, password)
}

func (t *Token) SetToken() (accessToken string, refreshToken string, err error) {
	return t.SetTokenFn()
}

func (t *Token) CheckToken(client *redis.Client, refreshPrefix string, request interface{}) (bool, error) {
	return t.CheckTokenFn(client, refreshPrefix, request)
}

func (t *Token) SaveToken(client *redis.Client, accessPrefix string, refreshPrefix string) error {
	return t.SaveTokenFn(client, accessPrefix, refreshPrefix)
}

func (t *Token) DeleteToken(client *redis.Client, accessPrefix string, refreshPrefix string) error {
	return t.DeleteTokenFn(client, accessPrefix, refreshPrefix)
}
