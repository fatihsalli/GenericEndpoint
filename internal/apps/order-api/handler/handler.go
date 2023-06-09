package handler

import (
	"GenericEndpoint/internal/apps/order-api"
	"GenericEndpoint/internal/models"
	"GenericEndpoint/pkg"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type Handler struct {
	MongoService   *order_api.MongoService
	ElasticService *order_api.ElasticService
}

func NewHandler(e *echo.Echo, mongoService *order_api.MongoService, elasticService *order_api.ElasticService) *Handler {
	router := e.Group("api/orders")
	h := &Handler{MongoService: mongoService, ElasticService: elasticService}

	//Routes
	router.GET("", h.GetAll)
	router.POST("", h.CreateOrder)
	router.POST("/GenericEndpoint", h.GenericEndpoint)
	router.POST("/GenericEndpointElastic", h.GenericEndpointElastic)
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
	orderList, err := h.MongoService.GetAll()

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
	filter, findOptions := h.MongoService.FromModelConvertToFilter(orderGetRequest)
	orderList, err := h.MongoService.GetOrdersWithFilter(filter, findOptions)

	if err != nil {
		c.Logger().Errorf("NotFoundError. %v", err.Error())
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("NotFoundError. %v", err.Error()),
		})
	}

	var orderResponse order_api.OrderResponse
	var orderResponseList []order_api.OrderResponse

	for _, order := range orderList {
		orderResponse.ID = order.ID
		orderResponse.UserID = order.UserID
		orderResponse.Status = order.Status
		orderResponse.City = order.City
		orderResponse.AddressDetail = order.AddressDetail
		orderResponse.Product = order.Product
		orderResponse.Total = order.Total

		if order.CreatedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
			orderResponse.CreatedAt = ""
		} else {
			orderResponse.CreatedAt = order.CreatedAt.String()
		}

		if order.UpdatedAt.String() == "0001-01-01 00:00:00 +0000 UTC" {
			orderResponse.UpdatedAt = ""
		} else {
			orderResponse.UpdatedAt = order.UpdatedAt.String()
		}

		orderResponseList = append(orderResponseList, orderResponse)
	}

	// Response success result data
	jsonSuccessResultData := models.JSONSuccessResultData{
		TotalItemCount: len(orderResponseList),
		Data:           orderResponseList,
	}

	c.Logger().Info("Orders are successfully listed.")
	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GenericEndpointElastic godoc
// @Summary get orders list with filter
// @ID get-orders-with-filter-from-elastic
// @Produce json
// @Param data body order_api.OrderGetRequest true "order filter data"
// @Success 200 {object} models.JSONSuccessResultData
// @Success 400 {object} pkg.BadRequestError
// @Success 404 {object} pkg.NotFoundError
// @Router /orders/GenericEndpointElastic [post]
func (h *Handler) GenericEndpointElastic(c echo.Context) error {
	var orderGetRequest order_api.OrderGetRequest

	if err := c.Bind(&orderGetRequest); err != nil {
		c.Logger().Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, pkg.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// Create filter and find options (exact filter,sort,field and match)
	orderList, err := h.ElasticService.GetFromElasticsearch(orderGetRequest)
	if err != nil {
		c.Logger().Errorf("InternalServerError. %v", err.Error())
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: fmt.Sprintf("InternalServerError. %v", err.Error()),
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

	// Create id and created date value
	orderModel.ID = uuid.New().String()
	orderModel.CreatedAt = time.Now()
	// We don't want to set null, so we put CreatedAt value.
	orderModel.UpdatedAt = orderModel.CreatedAt

	var total float64
	for _, product := range orderModel.Product {
		total = product.Price * float64(product.Quantity)
		orderModel.Total += total
	}

	result, err := h.MongoService.Insert(orderModel)

	if err != nil {
		c.Logger().Errorf("StatusInternalServerError: %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// Save to elasticsearch
	if err := h.ElasticService.SaveOrderToElasticsearch(orderModel); err != nil {
		c.Logger().Errorf("StatusInternalServerError (Elasticsearch) : %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong with elasticsearch!",
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

	result, err := h.MongoService.Delete(query)

	if err != nil || result == false {
		c.Logger().Errorf("NotFoundError. %v", err.Error())
		return c.JSON(http.StatusNotFound, pkg.NotFoundError{
			Message: fmt.Sprintf("NotFoundError. %v", err.Error()),
		})
	}

	// Delete from elasticsearch
	if err := h.ElasticService.DeleteOrderFromElasticsearch(query); err != nil {
		c.Logger().Errorf("StatusInternalServerError (Elasticsearch) : %v", err)
		return c.JSON(http.StatusInternalServerError, pkg.InternalServerError{
			Message: "Something went wrong with elasticsearch!",
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
