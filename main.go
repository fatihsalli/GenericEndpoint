package main

import (
	"GenericEndpoint/docs"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

// @title           Generic Get Endpoint
// @version         1.0
// @description     This is a generic get endpoint.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8010
// @BasePath  /api
func main() {
	e := echo.New()

	docs.SwaggerInfo.Host = "localhost:8010"
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	if err := e.Start(":8010"); err != nil {
		e.Logger.Fatal("Shutting down the server!")
	}
}
