package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// AuthMiddleware is a demo middleware function that checks for the presence of an Authorization header
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Unauthorized",
				})
			}
			// actual logic here (not implemented)
			return next(c)
		}
	}
}
