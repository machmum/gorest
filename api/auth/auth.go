package auth

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/machmum/gorest/utl/model"
	"golang.org/x/crypto/bcrypt"
	"runtime"
)

var (
	// return error in processing service
	ErrProfilePassword = errors.New("email and password did not match")
	ErrInvalidProfile  = errors.New("invalid profile")

	// redis attributes
	key, index string
)

func trace(c echo.Context) {
	if c != nil {
		function, file, line, _ := runtime.Caller(1)
		etrace := map[string]interface{}{
			"func": function,
			"file": file,
			"line": line,
		}
		c.Set("etrace", etrace)
	}
}

func (a *Auth) Login(c echo.Context, username string, password string) (result gorest.Profile, err error) {
	// leave error trace
	defer func() {
		if err != nil {
			trace(c)
		}
	}()

	prf, err := a.platform.postgresql.FindByUsername(a.conn.postgresql, username)
	if err != nil {
		return result, err
	}

	// validated id - password
	if prf.ID < 1 || len(prf.Password) < 1 {
		return result, ErrInvalidProfile
	}

	// Comparing the given password with the hash in db
	// bcrypt.CompareHashAndPassword({hash_db}, {given_pwd})
	err = bcrypt.CompareHashAndPassword([]byte(prf.Password), []byte(password))
	if err != nil {
		return result, ErrProfilePassword
	}

	result = gorest.Profile{
		ID:          prf.ID,
		FirstName:   prf.FirstName,
		LastName:    prf.LastName,
		Username:    prf.Username,
		Description: prf.Description,
		Email:       prf.Email,
	}

	// store to redis
	key = a.cfg.Prefix.Refresh + c.Get("oauth_refresh").(string)
	index = "profile"
	lifetime := a.cfg.Lifetime.Token * 2
	err = a.platform.rds.Save(a.conn.rds, result, key, index, lifetime)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (a *Auth) Logout(c echo.Context, token string) (err error) {
	// leave error trace
	defer func() {
		if err != nil {
			trace(c)
		}
	}()

	index = "profile"
	key = a.cfg.Prefix.Refresh + token
	err = a.platform.rds.Delete(a.conn.rds, key, index)
	if err != nil {
	}

	return err
}
