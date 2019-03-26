package api

import (
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

type APIValidator struct {
	validator *validator.Validate
}

func (av *APIValidator) Validate(i interface{}) error {
	return av.validator.Struct(i)
}

func (api *APIServer) router(web *echo.Echo) {
	web.Validator = &APIValidator{validator: validator.New()}

	v1 := web.Group("/api/v1")

	v1.POST("/users/register", api.PostRegister)
	v1.GET("/users/register/captcha", api.GetRegisterCaptcha)
}
