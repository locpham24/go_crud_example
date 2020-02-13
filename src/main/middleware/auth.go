package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"log"
	"main/model"
	"strconv"
)

func SetHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, "application/json")
		return next(c)
	}
}
func AuthWithJwt(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		claims := token.Claims.(jwt.MapClaims)
		userId, err := strconv.Atoi(claims["jti"].(string))

		if err != nil {

		}

		model.CurrentUser = model.User{
			Id:   userId,
			Name: claims["name"].(string),
		}

		log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

		return next(c)
	}
}
