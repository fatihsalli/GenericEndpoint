package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type OrderHandler struct {
	OrderCollection *mongo.Collection
}

func NewOrderHandler(e *echo.Echo, mongoCollection *mongo.Collection) *OrderHandler {
	router := e.Group("api/orders")
	h := &OrderHandler{OrderCollection: mongoCollection}

	//Routes
	router.GET("", h.GetOrders)

	return h
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	var order Order
	var orders []Order

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := h.OrderCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	for result.Next(ctx) {
		if err := result.Decode(&order); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}
