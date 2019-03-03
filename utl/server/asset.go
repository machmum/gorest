package server

import (
	"image"
	"net/http"

	// image decoder
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func GetImage(source string) (int, int, error) {
	resp, err := http.Get(source)
	if err != nil {
		return 0, 0, err
	}

	// decode image
	m, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	g := m.Bounds()

	return g.Dx(), g.Dy(), nil
}
