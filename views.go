package api

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"git.iptq.io/nso/common/models"
	"github.com/dchest/captcha"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInfo struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=6"`
	Captcha  string `json:"captcha" validate:"required"`
}

func (api *APIServer) PostRegister(c echo.Context) (err error) {
	// get the session
	sess, err := session.Get("session", c)
	if err != nil {
		return
	}

	// validate form info
	info := new(RegisterInfo)
	if err = c.Bind(info); err != nil {
		return
	}
	if err = c.Validate(info); err != nil {
		return
	}

	// validate the captcha
	captchaId := sess.Values["captcha"]
	if !captcha.VerifyString(captchaId, info.Captcha) {
		return fmt.Errorf("Invalid captcha")
	}

	email := strings.ToLower(info.Email)
	username := strings.ToLower(info.Username)
	password := []byte(info.Password)

	// does the user already exist?
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
		ID uint `json:"id"`
	}{
		ID: user.ID,
	})
}

func (api *APIServer) GetRegisterCaptcha(c echo.Context) (err error) {
	// get the session
	sess, err := session.Get("session", c)
	if err != nil {
		return
	}

	// generate captcha
	id := captcha.New()
	sess.Values["captcha"] = id
	sess.Save(c.Request(), c.Response())

	// return image
	r, w := io.Pipe()
	go func() {
		captcha.WriteImage(w, id, captcha.StdWidth, captcha.StdHeight)
		w.Close()
	}()
	return c.Stream(http.StatusOK, "image/png", r)
}
