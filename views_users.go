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

type UserInfo struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	UsernameCase string `json:"usernameCase"`
}

func GetUserInfo(user *models.User) UserInfo {
	return UserInfo{
		ID:           user.ID,
		Username:     user.Username,
		UsernameCase: user.UsernameCase,
	}
}

type LoginInfo struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

func (api *APIServer) PostLogin(c echo.Context) (err error) {
	// get the session
	sess, err := session.Get("session", c)
	if err != nil {
		return
	}

	// validate form info
	info := new(LoginInfo)
	if err = c.Bind(info); err != nil {
		return
	}
	if err = c.Validate(info); err != nil {
		return
	}

	identifier := strings.ToLower(info.Identifier)
	password := []byte(info.Password)

	// get the user if exists
	var user models.User
	api.db.Where("username = ?", identifier).Or("email = ?", identifier).First(&user)

	// check the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), password); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Incorrect credentials.")
	}

	// now put it in the session!
	sess.Values["user_id"] = user.ID
	sess.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, GetUserInfo(&user))
}

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
	captchaId, ok := sess.Values["captcha"]
	if !ok || !captcha.VerifyString(captchaId.(string), info.Captcha) {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid captcha")
	}

	email := strings.ToLower(info.Email)
	usernameCase := info.Username
	username := strings.ToLower(info.Username)
	password := []byte(info.Password)

	// does the user already exist?
	var count uint
	api.db.Where("username = ?", username).Or("email = ?", email).Find(&models.User{}).Count(&count)
	if count > 0 {
		return echo.NewHTTPError(http.StatusConflict, "A user with this username or email exists already.")
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	// god damn it osu
	osuDumbHash := fmt.Sprintf("%x", md5.Sum(password))
	osuHash, err := bcrypt.GenerateFromPassword([]byte(osuDumbHash), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user := models.User{
		Email:        email,
		UsernameCase: usernameCase,
		Username:     username,
		Password:     string(hash),
		OsuPassword:  string(osuHash),
	}
	api.db.Create(&user)

	return c.JSON(http.StatusOK, GetUserInfo(&user))
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
