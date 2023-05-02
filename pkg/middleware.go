package pkg

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			switch err.(type) {
			case *NotFoundError:
				return c.JSON(http.StatusNotFound, NotFoundError{
					Message: fmt.Sprintf("NotFoundError: %v", err.Error()),
				})
			case *BadRequestError:
				return c.JSON(http.StatusBadRequest, BadRequestError{
					Message: fmt.Sprintf("BadRequestError: %v", err.Error()),
				})
			case *ClientSideError:
				return c.JSON(http.StatusBadRequest, ClientSideError{
					Message: fmt.Sprintf("ClientSideError: %v", err.Error()),
				})
			default:
				return c.JSON(http.StatusInternalServerError, InternalServerError{
					Message: fmt.Sprintf("StatusInternalServerError: %v", err.Error()),
				})
			}
		}
		return nil
	}
}
