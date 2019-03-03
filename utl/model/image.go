package gorest

type (
	// size request
	SizeReq struct {
		Menu     *WHReq `json:"menu,omitempty"`
		Profile  *WHReq `json:"profile,omitempty"`
		Product  *WHReq `json:"product,omitempty"`
		Category *WHReq `json:"category,omitempty"`
	}

	// width x height
	WHReq struct {
		Width  int `json:"width" validate:"numeric"`
		Height int `json:"height" validate:"numeric"`
	}

	// image response
	ImageRes struct {
		URL  string `json:"url"`
		Size *WHReq `json:"size"`
	}

	// image product
	ImageProduct struct {
		Preview   *ImageRes `json:"preview,omitempty"`
		Thumbnail *ImageRes `json:"thumbnail,omitempty"`
	}

	ImageProductSlice []ImageProduct
)
