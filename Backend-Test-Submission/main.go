package main

import (
	"Backend-Test-Submission/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitDB()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Hello Client!",
		})
	})
	e.Logger.Fatal(e.Start(":8080"))
}
