package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	address, ok := os.LookupEnv("HDRCRYPTO_SERVER_ADDR")
	if !ok {
		address = ":3000"
	}
	e.Logger.Fatal(e.Start(address))
}
