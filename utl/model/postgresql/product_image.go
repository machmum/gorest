package gorestdb

type (
	ProductImage struct {
		// Product Image table
		ID        uint   `json:"id" gorm:"primary_key"`
		ProductID uint   `json:"product_id" gorm:"column:product_id"`
		Name      string `json:"name"`
		Status    int    `json:"status"`
		Base
	}
)
