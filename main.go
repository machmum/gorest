package main

import (
	"github.com/machmum/gorest/api"
	"github.com/machmum/gorest/config"
)

func main() {

	path := "./config.local.yaml"

	// get config
	cfg, err := config.Load(path)
	if err != nil {
		panic(err.Error())
	}

	// run engine
	api.Start(cfg)
}
