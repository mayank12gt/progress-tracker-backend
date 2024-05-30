package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	database "github.com/mayank12gt/ProgressTracker/db"
	"github.com/mayank12gt/ProgressTracker/model"
	"gorm.io/gorm"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("Logger Hit")
		return next(c)
	}
}

func VerifyJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("auth middleware")

		//println("auth middleware")
		cookie, err := c.Cookie("jwt")

		if err != nil {
			return c.String(400, "unauthenticated")
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("SecretKey"), nil
		})

		if err != nil {

			return c.String(400, "unauthenticated")
		}

		claims := token.Claims.(*jwt.StandardClaims)

		var user model.User

		database.Database.Db.Session(&gorm.Session{}).Where("id = ?", claims.Issuer).First(&user)
		c.Set("user", user)
		return next(c)

	}

}
