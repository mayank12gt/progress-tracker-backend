package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	database "github.com/mayank12gt/ProgressTracker/db"
	"github.com/mayank12gt/ProgressTracker/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDTO struct {
	ID uint `json:"id"`

	Name string `json:"name"`

	Email string `json:"email"`

	Password string `json:"password"`
}

type SignInReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func createUserDTO(user model.User) UserDTO {

	return UserDTO{ID: user.ID, Name: user.Name, Email: user.Email, Password: user.Password}
}

func FindUser(email string) model.User {
	var user model.User
	database.Database.Db.Session(&gorm.Session{}).Where("email=?", email).Find(&user)
	// database.Database.Db.Session(&gorm.Session{}).First(&user, email)
	return user
}

func SignUp(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	user.Password = string(password)
	database.Database.Db.Create(&user)
	userResponse := createUserDTO(user)
	return c.JSON(http.StatusOK, userResponse)
}

func SignIn(c echo.Context) error {
	var signInReq SignInReq

	if err := c.Bind(&signInReq); err != nil {
		return c.String(400, err.Error())
	}
	var user model.User
	user = FindUser(string(signInReq.Email))
	if user.ID == 0 {
		return c.String(400, "User Not Found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInReq.Password)); err != nil {

		return c.String(400, "Incorrect Password")
	}
	//userRes := CreateUserDTO(user)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte("SecretKey"))

	if err != nil {

		return c.String(400, "Unable to log In")
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Secure:   false,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Path:     "/",
	} //Creates the cookie to be passed.

	c.SetCookie(&cookie)

	return c.String(200, "login success")
}

func SignOut(c echo.Context) error {

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), //Sets the expiry time an hour ago in the past.
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
	}

	c.SetCookie(&cookie)

	return c.String(200, "loggedOut")
}

func Profile(c echo.Context) error {
	println("profile")
	user := c.Get("user")
	//userRes := createUserDTO(model.User(user))

	return c.JSON(200, user)

}
