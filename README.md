# Gorest API

Golang starter kit for develop RESTful API based on [Gorsk](https://github.com/ribice/gorsk).
The initial process using **OAuth 2.0 Password Grant**. 

Depedency used :
```bash
|-------------------------------------|--------------------------------------------|--------------|
|             DEPENDENCY              |                  REPOURL                   |   LICENSE    |
|-------------------------------------|--------------------------------------------|--------------|
| github.com/labstack/echo            | https://github.com/labstack/echo           | MIT          |
| golang.org/x/crypto/bcrypt          | https://github.com/golang/crypto           |              |
| gopkg.in/go-playground/validator.v8 | https://github.com/go-playground/validator | MIT          |
| github.com/stretchr/testify         | https://github.com/stretchr/testify        | MIT          |
| github.com/go-redis/redis           | https://github.com/go-redis/redis          | Other        |
| github.com/alicebob/miniredis       | https://github.com/alicebob/miniredis      | MIT        |
| github.com/gofrs/uuid               | https://github.com/gofrs/uuid              | Other        |
| github.com/jinzhu/gorm              | https://github.com/jinzhu/gorm             | MIT          |
| github.com/sirupsen/logrus          | https://github.com/sirupsen/logrus         | MIT          |
| gopkg.in/yaml.v2                    | https://github.com/go-yaml/yaml            |              |
| go.uber.org/zap                     | https://github.com/uber-go/zap             | Other        |
|-------------------------------------|--------------------------------------------|--------------|
```

1. Echo - HTTP 'framework'.
2. Bcrypt - Password hashing
3. Validator - Request validation
4. Testify/Assert - Asserting test results
5. Go-redis - Type-safe Redis client for Golang
6. Miniredis - Pure Go Redis server for Go unittests
7. Gofrs/uuid - pure Go implementation of Universally Unique Identifiers (UUID)
8. Jinzhu/gorm - ORM for golang
9. Sirupsen/logrus - Structured, pluggable logging for Go
10. Yaml - Unmarshalling YAML config file
11. Zap - Blazing fast, structured, leveled logging in Go

## Getting started

Using Gorest requires having Go 1.10 or above.

1. To use Gorest as a starting point of a real project whose package name is something like `github.com/author/project`, move the directory `$GOPATH/github.com/machmum/gorest` to `$GOPATH/github.com/author/project` and do a global replacement of the string `github.com/machmum/gorest` with `github.com/author/project`.

2. Rename the gorest package inside `utl/model` with your own project name, then using search & replace do a global replacement of `.gorest` with your project name.

3. Change the configuration file `config.local.yaml` according to your needs, or create a new one.

4. In `cmd/migration/main.go` set up psn variable and then run it (go run main.go). It will create all tables, and necessery data. Use `admin/admin` to get token and `username/user` to login.

5. Make sure `dep` already installed, and install package above using `dep ensure -v`

6. Run the app using:
```bash
go run main.go
```

The application runs as an HTTP server at port 8080. It provides the following RESTful endpoints:

* `POST /v1/oauth/token/request`: accepts `{"username":"admin","password":"admin"}` and returns access_token, refresh_token, token_type and expiry_access.

* `POST /v1/oauth/token/refresh`: accepts `{"username":"admin","password":"admin","refresh_token":"token"}`, refresh_token `token` got from previous endpoint. Returns access_token, refresh_token, token_type and expiry_access.

* `POST /v1/auth/login`: accepts `{"username":"username","password":"user","scope":"profile"}`. Returns profile data.

* `GET /v1/auth/logout`: returns info logged out.

* `POST /v1/profile`: accepts `{"category":"user","scope":"profile"}`. Returns profile data.

* `POST /v1/profile/detail`: accepts `{"scope":"profile+product","profile_id":1,"product_id":1,"size":{"profile":{"width":100,"height":0},"product":{"width":400,"height":0}}}`. Returns profile data.

* `POST /v1/product`: accepts `{"scope":"product","profile_id":1,"product_id":1,"size":{"product":{"width":400,"height":0}}}`. Returns product data.

All endpoint accepts Headers `Content-type: application/json`, `Authorization: Bearer token`, except oauth endpoint, remove Headers Authorization. `token` is `access_token` got from `/v1/oauth/token/request`




















