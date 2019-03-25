package api

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type APIServer struct {
	config *Config
	web    *echo.Echo
}

func NewInstance(config *Config) (api *APIServer, err error) {
	web := echo.New()
	router(web)

	// logging middleware
	web.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	api = &APIServer{
		config: config,
		web:    web,
	}
	return
}

func (api *APIServer) Run() {
	log.Fatal(api.web.Start(api.config.BindAddr))
}
