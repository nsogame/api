package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	"git.iptq.io/nso/common/models"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInfo struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (api *APIServer) PostRegister(c echo.Context) (err error) {
	// validate form info
	info := new(RegisterInfo)
	if err = c.Bind(info); err != nil {
		return
	}
	if err = c.Validate(info); err != nil {
		return
	}

	email := strings.ToLower(info.Email)
	username := strings.ToLower(info.Username)
	password := []byte(info.Password)

	// does the user exist?
	var count uint
	api.db.Where("username = ?", username).Or("email = ?", email).Find(&models.User{}).Count(&count)
	if count > 0 {
		return fmt.Errorf("A user with this username or email exists already.")
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	// god damn it osu
	osuDumbHash := md5.Sum(password)
	osuHash, err := bcrypt.GenerateFromPassword(osuDumbHash[:], bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user := &models.User{
		Email:       email,
		Username:    username,
		Password:    string(hash),
		OsuPassword: string(osuHash),
	}
	api.db.Create(&user)

	return c.JSON(http.StatusOK, struct {
		ID int
	}{
		ID: 5,
	})
}
