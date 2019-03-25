package api

import (
	"git.iptq.io/nso/api/views"
	"github.com/labstack/echo"
)

func router(web *echo.Echo) {
	web.POST("/api/users/register", views.PostRegister)
}
