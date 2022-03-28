package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hdrcrypto/cmd"
	"net/http"
	"os"
)

func main() {
	cmd.Execute()
}

func serve() {
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
