package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func BodyLogger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		requestID := c.Response().Header().Get(echo.HeaderXRequestID)
		userUUID, _ := c.Get(UserUUID).(string)

		log := logger.Info().
			Bytes("request_body", reqBody).
			Bytes("response_body", resBody)

		if requestID != "" {
			log.Str("request_id", requestID)
		}

		if userUUID != "" {
			log.Str(UserUUID, userUUID)
		}

		log.Send()
	})
}
