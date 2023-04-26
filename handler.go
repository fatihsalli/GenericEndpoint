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
	router.POST("", h.GetOrders)

	return h
}

// GetOrders godoc
// @Summary get order list with filter
// @ID get-orders
// @Produce json
// @Param data body OrderRequest true "order filter"
// @Success 200 {array} Order
// @Success 400
// @Success 500
// @Router /orders [post]
func (h *OrderHandler) GetOrders(c echo.Context) error {
	var orderRequest OrderRequest
	if err := c.Bind(&orderRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding: %v", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := h.OrderCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var order Order
	var orders []Order

	for result.Next(ctx) {
		if err := result.Decode(&order); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}
