package route

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"main/api"
	"main/conf"
	myMw "main/middleware"
)

func Init() *echo.Echo {
	e := echo.New()
	taskGroup := e.Group("/task")

	taskGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    conf.JwtKey,
	}))

	taskGroup.Use(myMw.AuthWithJwt)
	taskGroup.Use(myMw.SetHeader)

	taskGroup.GET("", api.GetAll())
	taskGroup.POST("", api.Create())
	taskGroup.GET("/:id", api.GetOne())
	taskGroup.PUT("/:id", api.Update())
	taskGroup.DELETE("/:id", api.Delete())

	e.POST("/login", api.Login())

	return e
}
