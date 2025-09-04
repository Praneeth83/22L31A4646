package routes

import "github.com/labstack/echo/v4"

func ShorterRoutes(e *echo.Echo) {
	e.POST("/shorturls", controllers.GetShortUrl)
	e.GET("/shorturls/:code", controllers.GetStats)
	e.GET("/:code", controllers.UrlRedirecter)
}
