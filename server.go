package main

import (
	"log"
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/crestron-control-microservice/crestroncontrol"
	"github.com/byuoitav/crestron-control-microservice/handlers"
	"github.com/byuoitav/hateoas"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	var err error

	crestroncontrol.SignalConfigFile, err = crestroncontrol.ParseConfig()
	if err != nil {
		log.Fatal("Error parsing signal-configuration.json: " + err.Error())
	}

	// setting up a web server
	port := ":8004"

	//the router determines which code gets executed after a request
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	secure.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	secure.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	secure.GET("/:address/power/on", handlers.PowerOn)
	secure.GET("/:address/power/standby", handlers.Standby)
	secure.GET("/:address/input/:port", handlers.SwitchInput)
	secure.GET("/:address/volume/set/:value", handlers.SetVolume)
	secure.GET("/:address/volume/mute", handlers.VolumeMute)
	secure.GET("/:address/volume/unmute", handlers.VolumeUnmute)
	secure.GET("/:address/display/blank", handlers.BlankDisplay)
	secure.GET("/:address/display/unblank", handlers.UnblankDisplay)
	secure.GET("/:address/volume/get", handlers.GetVolume)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
