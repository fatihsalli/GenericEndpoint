package handler

import (
	"GenericEndpoint/internal/apps/order-api"
	"GenericEndpoint/internal/models"
	"GenericEndpoint/pkg"
	"fmt"
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
	router.POST("/GenericEndpoint", h.GenericEndpoint)
	router.DELETE("/:id", h.DeleteOrder)

	return h
}

// GetAll godoc
// @Summary get all order list
// @ID get-all
// @Produce json
// @Success 200 {object} models.JSONSuccessResultData
// @Success 500 {object} pkg.InternalServerError
// @Router /orders [get]
func (h *Handler) GetAll(c echo.Context) error {
	orderList, err := h.Service.GetAll()

	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(orderList),
		Data:           orderList,
	}

	c.Logger().Info("All orders are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GenericEndpoint godoc
// @Summary get orders list with filter
// @ID get-orders-with-filter
// @Produce json
// @Param data body order_api.OrderGetRequest true "order filter data"
// @Success 200 {object} models.JSONSuccessResultData
// @Success 400 {object} pkg.BadRequestError
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/GenericEndpoint [post]
func (h *Handler) GenericEndpoint(c echo.Context) error {
	var orderGetRequest order_api.OrderGetRequest

	if err := c.Bind(&orderGetRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Create filter and find options (exact filter,sort,field and match)
	filter, findOptions := h.Service.FromModelConvertToFilter(orderGetRequest)

	orderList, err := h.Service.GetOrdersWithFilter(filter, findOptions)

	if err != nil {
		c.Logger().Errorf("NotFoundError. %v", err.Error())
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("NotFoundError. %v", err.Error()),
		})
	}

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(orderList),
		Data:           orderList,
	}

	c.Logger().Info("Orders are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// CreateOrder godoc
// @Summary add a new item to the order list
// @ID create-order
// @Produce json
// @Param data body order_api.OrderCreateRequest true "order data"
// @Success 201 {object} models.JSONSuccessResultId
// @Success 400 {object} pkg.BadRequestError
// @Success 500 {object} pkg.InternalServerError
// @Router /orders [post]
func (h *Handler) CreateOrder(c echo.Context) error {
	var orderCreateRequest order_api.OrderCreateRequest

	if err := c.Bind(&orderCreateRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	var orderModel models.Order

	orderModel.UserID = orderCreateRequest.UserID
	orderModel.Status = orderCreateRequest.Status
	orderModel.City = orderCreateRequest.City
	orderModel.AddressDetail = orderCreateRequest.AddressDetail
	orderModel.Product = orderCreateRequest.Product

	result, err := h.Service.Insert(orderModel)

	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	c.Logger().Infof("{%v} with id is created.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
}

// DeleteOrder godoc
// @Summary delete an order item by ID
// @ID delete-order-by-id
// @Produce json
// @Param id path string true "order ID"
// @Success 200 {object} models.JSONSuccessResultId
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(c echo.Context) error {
	query := c.Param("id")

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		c.Logger().Errorf("NotFoundError. %v", err.Error())
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("NotFoundError. %v", err.Error()),
		})
	}

	// To response id and success boolean
	jsonSuccessResultId := models.JSONSuccessResultId{
		ID:      query,
		Success: true,
	}

	c.Logger().Infof("{%v} with id is deleted.", jsonSuccessResultId.ID)
	return c.JSON(http.StatusCreated, jsonSuccessResultId)
}
