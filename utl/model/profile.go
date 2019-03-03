package gorest

type (
	// Struct represent /category request
	ProfileReq struct {
		Category string `json:"category" validate:"required"`
		Scope    string `json:"scope,required" validate:"required"`
	}

	PflDetailReq struct {
		Scope     string   `json:"scope" validate:"required"`
		ProfileID int      `json:"profile_id" validate:"required"`
		ProductID int      `json:"product_id"`
		Size      *SizeReq `json:"size"`
	}

	// Struct represent response
	Profile struct {
		ID          uint   `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Username    string `json:"username"`
		Description string `json:"description,omitempty"`
		Email       string `json:"email,omitempty"`

		// Product detail
		Products ProductSlice `json:"products,omitempty"`
	}

	ProfileSlice []Profile
)
