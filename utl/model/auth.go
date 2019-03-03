package gorest

type (
	// struct represent auth/login request
	LoginReq struct {
		Username string
		Password string
		Scope    string
	}
)
