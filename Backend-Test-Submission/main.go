package main

import (
	"Backend-Test-Submission/config"
	"Backend-Test-Submission/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	config.InitDB()
	e := echo.New()
	routes.ShorterRoutes(e)
	e.Logger.Fatal(e.Start(":8080"))
}
