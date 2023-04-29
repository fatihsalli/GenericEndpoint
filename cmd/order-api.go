package cmd

import (
	"GenericEndpoint/internal/apps/order-api"
	"GenericEndpoint/internal/apps/order-api/handler"
	"GenericEndpoint/internal/configs"
	"GenericEndpoint/internal/repository"
	"GenericEndpoint/pkg"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

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

	// Start server as asynchronous
	go func() {
		if err := e.Start(config.Server.Port["orderAPI"]); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server!")
		}
	}()

	// Graceful Shutdown
	pkg.GracefulShutdown(e, 10*time.Second)
}
