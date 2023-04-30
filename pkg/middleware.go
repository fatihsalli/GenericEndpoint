package pkg

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func ErrorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			switch e := err.(type) {
			case *NotFoundError:
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": e.Error(),
				})
			case *BadRequestError:
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": e.Error(),
				})
			case *ClientSideError:
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": e.Error(),
				})
			default:
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"error": "Internal server error",
				})
			}
		}
		return nil
	}
}
