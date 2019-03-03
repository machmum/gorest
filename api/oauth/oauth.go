package oauth

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/server"
	"runtime"
)

const (
	// token grant_type
	GrantAccess  = "request"
	GrantRefresh = "refresh"
	TokenType    = "Bearer"
)

// set error trace
func trace(c echo.Context) {
	if c != nil {
		function, file, line, _ := runtime.Caller(2)
		etrace := map[string]interface{}{
			"func": function,
			"file": file,
			"line": line,
		}
		c.Set("etrace", etrace)
	}
}

func (o *Oauth) Tokenize(c echo.Context, gtp, username, password, rtk string) (result *gorest.OauthToken, status bool, err error) {
	var (
		channel    *gorestdb.OauthChannel   // init nil pointer to struct
		oauthToken = new(gorest.OauthToken) // create new pointer to empty struct

		// token redis
		key, index    string
		tac, trf, ttp string
		exp           uint
	)

	// leave error trace
	defer func() {
		if err != nil {
			trace(c)
		}
	}()

	// get channel in redis
	index = "channel_" + username
	key = o.cfg.Prefix.Apps

	if o.platform.rds != nil {
		rds, err := o.platform.rds.Find(o.conn.rds, key, index)
		if rds != "" {

			err = json.Unmarshal([]byte(rds), &channel)
			if err != nil {
				return nil, false, err
			}

			if channel.Password != password {
				return nil, false, server.ErrInvalidUsernamePassword
			}
		}
	}

	// not found channel in redis
	// get channel from database
	if channel == nil {

		channel, err = o.platform.postgresql.FindByUsername(o.conn.postgresql, username)
		if err != nil {
			return nil, false, err
		}

		if channel.Password != password {
			return nil, false, server.ErrInvalidUsernamePassword
		}

		err = o.platform.rds.Save(o.conn.rds, channel, key, index, o.cfg.Lifetime.Apps)
		if err != nil {
			return nil, false, err
		}
	}

	// initiate token
	// with given lifetime
	ttp, exp = o.token.InitToken(o.cfg.Lifetime.Token, TokenType, channel.Username, channel.Password)

	// get access token
	if gtp == GrantAccess {
		if tac, trf, err = o.token.SetToken(); err != nil {
			return nil, false, err
		}
	}

	// get refresh token
	if gtp == GrantRefresh {
		request := map[string]interface{}{
			"username":      username,
			"password":      password,
			"refresh_token": rtk,
		}

		// check token
		// set token.SavedAccess and token.SavedRefresh
		if status, err = o.token.CheckToken(o.conn.rds, o.cfg.Prefix.Refresh, request); err != nil {
			return nil, status, err
		}

		// set new token
		if tac, trf, err = o.token.SetToken(); err != nil {
			return nil, false, err
		}

		// delete token
		if err = o.token.DeleteToken(o.conn.rds, o.cfg.Prefix.Access, o.cfg.Prefix.Refresh); err != nil {
			return nil, false, err
		}
	}

	// save token to redis
	if err = o.token.SaveToken(o.conn.rds, o.cfg.Prefix.Access, o.cfg.Prefix.Refresh); err != nil {
		return nil, false, err
	}

	oauthToken = &gorest.OauthToken{
		AccessToken:  tac,
		RefreshToken: trf,
		TokenType:    ttp,
		ExpireAccess: exp,
	}

	// prepare result
	return oauthToken, false, nil
}
