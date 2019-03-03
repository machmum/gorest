package gorestdb

type (
	// Product table
	Product struct {
		Base
		ID          uint   `json:"id" gorm:"primary_key"`
		ProfileID   uint   `json:"profile_id" gorm:"column:profile_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		NormalPrice string `json:"normal_price"`
		SalePrice   string `json:"sale_price"`
		Status      int    `json:"status"`
	}

	ProductDetail struct {
		ProductID              uint `json:"product_id" gorm:"column:product_id"`
		ProductName            string
		NormalPrice, SalePrice string
		ProductImage           string
	}
)
