package gorest

type (
	// struct represent oauth/token request
	Credentials struct {
		Username     string `json:"username" validate:"required"`
		Password     string `json:"password" validate:"required"`
	}

	CredentialsRefresh struct {
		Username     string `json:"username" validate:"required"`
		Password     string `json:"password" validate:"required"`
		RefreshToken string `json:"refresh_token" validate:"omitempty"`
	}

	// struct represent oauth/token response
	OauthToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpireAccess uint   `json:"expiry_access"`
	}
)
