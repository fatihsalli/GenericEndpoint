package handler

import (
	"GenericEndpoint/models"
	"context"
	"github.com/google/uuid"
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
	router.POST("/get", h.GetOrders)
	router.POST("", h.CreateOrder)

	return h
}

// GetOrders godoc
// @Summary get order list with filter
// @ID get-orders
// @Produce json
// @Param data body models.OrderRequest true "order filter"
// @Success 200 {array} models.Order
// @Success 400
// @Success 500
// @Router /orders/get [post]
func (h *OrderHandler) GetOrders(c echo.Context) error {
	var orderRequest models.OrderRequest
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

	var order models.Order
	var orders []models.Order

	for result.Next(ctx) {
		if err := result.Decode(&order); err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		orders = append(orders, order)
	}

	return c.JSON(http.StatusOK, orders)
}

// CreateOrder godoc
// @Summary add a new item to the order list
// @ID create-order
// @Produce json
// @Param data body models.OrderCreateRequest true "order data"
// @Success 201
// @Success 400
// @Success 500
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var orderCreateRequest models.OrderCreateRequest

	if err := c.Bind(&orderCreateRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var orderModel models.Order

	orderModel.ID = uuid.New().String()
	orderModel.UserID = orderCreateRequest.UserID
	orderModel.Status = orderCreateRequest.Status
	orderModel.City = orderCreateRequest.City
	orderModel.AddressDetail = orderCreateRequest.AddressDetail
	orderModel.Product = orderCreateRequest.Product
	orderModel.CreatedAt = time.Now()
	orderModel.UpdatedAt = orderModel.CreatedAt

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := h.OrderCollection.InsertOne(ctx, orderModel)

	if result.InsertedID == nil || err != nil {
		c.Logger().Errorf("Bad Request. Something went wrong", err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orderModel.ID)
}
