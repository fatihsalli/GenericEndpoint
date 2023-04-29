package cmd

import (
	"GenericEndpoint/docs"
	"GenericEndpoint/internal/apps/order-api"
	"GenericEndpoint/internal/apps/order-api/handler"
	"GenericEndpoint/internal/configs"
	"GenericEndpoint/internal/repository"
	"GenericEndpoint/pkg"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"time"
)

// @title           Echo Order API
// @version         1.0
// @description     This is an order API for generic query.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8011
// @BasePath  /api
func StartOrderAPI() {
	// Echo instance
	e := echo.New()

	// Get config
	config := configs.GetConfig("test")

	// Create repo and service
	mongoOrderCollection := configs.ConnectDB(config.Database.Connection).Database(config.Database.DatabaseName).Collection(config.Database.OrderCollectionName)
	OrderRepository := repository.NewRepository(mongoOrderCollection)
	OrderService := order_api.NewService(OrderRepository)

	// Create handler
	handler.NewHandler(e, OrderService)

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8011"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	// Start server as asynchronous
	go func() {
		if err := e.Start(config.Server.Port["orderAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	// Graceful Shutdown
	pkg.GracefulShutdown(e, 10*time.Second)
}
