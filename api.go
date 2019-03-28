package api

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
	"github.com/nsogame/common"
)

type APIServer struct {
	config *Config
	db     *common.DB
	web    *echo.Echo
}

func NewInstance(config *Config) (api *APIServer, err error) {
	// db
	db, err := common.ConnectDB(config.DbProvider, config.DbConnection)
	if err != nil {
		return
	}

	web := echo.New()
	web.Debug = config.Debug

	// middleware
	web.Use(session.Middleware(sessions.NewCookieStore([]byte(config.SecretKey))))
	web.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	web.Use(middleware.Recover())
	web.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1234", "https://osu.ppy.sh"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	api = &APIServer{
		config: config,
		db:     db,
		web:    web,
	}
	api.router(web)
	return
}

func (api *APIServer) Close() {
	api.db.Close()
}

func (api *APIServer) Run() {
	defer api.Close()
	log.Fatal(api.web.Start(api.config.BindAddr))
}
