package gorest

type (
	// Struct represent response
	Product struct {
		ID     uint              `json:"id"`
		Name   string            `json:"name"`
		Price  *Price            `json:"price"`
		Image  *ImageProduct     `json:"image,omitempty"`
		Images ImageProductSlice `json:"images,omitempty"`
	}

	ProductSlice []Product

	Price struct {
		Normal string `json:"normal,omitempty"`
		Sale   string `json:"sale,omitempty"`
	}
)
