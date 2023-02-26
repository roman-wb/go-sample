package middleware

import "github.com/labstack/echo/v4"

const (
	UserUUIDHeader = "User-Uuid"
	UserUUID       = "user_uuid"
)

func ExtractUserUUID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header[UserUUIDHeader]
			if len(header) > 0 {
				c.Set(UserUUID, header[0])
			}

			return next(c)
		}
	}
}
