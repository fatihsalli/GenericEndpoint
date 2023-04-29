package handler

import (
	"GenericEndpoint/internal/apps/order-api"
	"GenericEndpoint/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	Service *order_api.Service
}

func NewHandler(e *echo.Echo, service *order_api.Service) *Handler {
	router := e.Group("api/orders")
	h := &Handler{Service: service}

	//Routes
	router.GET("", h.GetAll)
	router.POST("", h.CreateOrder)

	return h
}

// GetAll godoc
// @Summary get all order list
// @ID get-all
// @Produce json
// @Success 200 {array} models.Order
// @Success 500
// @Router /orders [get]
func (h *Handler) GetAll(c echo.Context) error {
	orderList, err := h.Service.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, orderList)
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
func (h *Handler) CreateOrder(c echo.Context) error {
	var orderCreateRequest order_api.OrderCreateRequest

	if err := c.Bind(&orderCreateRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var orderModel models.Order

	orderModel.UserID = orderCreateRequest.UserID
	orderModel.Status = orderCreateRequest.Status
	orderModel.City = orderCreateRequest.City
	orderModel.AddressDetail = orderCreateRequest.AddressDetail
	orderModel.Product = orderCreateRequest.Product

	result, err := h.Service.Insert(orderModel)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result.ID)
}
