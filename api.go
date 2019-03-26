package api

import (
	"log"
	"net/http"

	"git.iptq.io/nso/common/models"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
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

	// middleware
	web.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SecretKey))))
	web.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	web.Use(middleware.Recover())
	web.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:1234"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
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
