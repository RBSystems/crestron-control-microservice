package main

import (
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/hateoas"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	port := ":8004"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	secure.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	secure.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	// secure.GET("/:address/power/on", handlers.PowerOn)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
