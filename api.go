package api

import (
	"log"

	"git.iptq.io/nso/common/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type APIServer struct {
	config *Config
	db     *gorm.DB
	web    *echo.Echo
}

func NewInstance(config *Config) (api *APIServer, err error) {
	db, err := gorm.Open(config.DbProvider, config.DbConnection)
	if err != nil {
		return
	}

	// TODO: remvoe later
	db.AutoMigrate(&models.User{})

	web := echo.New()
	web.Debug = config.Debug

	// logging middleware
	web.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	api = &APIServer{
		config: config,
		db:     db,
		web:    web,
	}
	api.router(web)
	return
}

func (api *APIServer) Run() {
	log.Fatal(api.web.Start(api.config.BindAddr))
}
