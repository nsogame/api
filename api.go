package api

import (
	"log"

	"github.com/labstack/echo"
)

type APIServer struct {
	config *Config
	web    *echo.Echo
}

func NewInstance(config *Config) (api *APIServer, err error) {
	web := echo.New()
	router(web)

	api = &APIServer{
		config: config,
		web:    web,
	}
	return
}

func (api *APIServer) Run() {
	log.Fatal(api.web.Start(api.config.BindAddr))
}
