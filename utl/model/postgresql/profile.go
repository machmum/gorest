package gorestdb

type (
	Profile struct {
		// Profile table
		ID          uint   `json:"id" gorm:"primary_key"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Password    string `json:"password,omitempty" validate:"required"`
		Username    string `json:"username" validate:"required"`
		Description string `json:"description"`
		Email       string `json:"email"`
		Category    string `json:"category"`
		Base
	}

	ProfileSimple struct {
		ID                  uint `json:"id"`
		FirstName, LastName string
		Username            string
		ProductID           uint   `json:"pid,omitempty" gorm:"column:product_id"`
		ProductName         string `json:",omitempty"`
		NormalPrice         string `json:",omitempty"`
		SalePrice           string `json:",omitempty"`
		ProductImageID      uint   `json:"pi_id,omitempty" gorm:"column:pi_id"`
		ProductImage        string `json:",omitempty"`
	}

	ProfileSimpleSlice []ProfileSimple
)
