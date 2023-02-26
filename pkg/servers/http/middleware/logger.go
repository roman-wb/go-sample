package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func Logger(logger *zerolog.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogProtocol:      true,
		LogRemoteIP:      true,
		LogHost:          true,
		LogMethod:        true,
		LogURIPath:       true,
		LogRoutePath:     true,
		LogRequestID:     true,
		LogReferer:       true,
		LogUserAgent:     true,
		LogStatus:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogError:         true,
		LogValuesFunc: func(c echo.Context, vals middleware.RequestLoggerValues) error {
			userUUID, _ := c.Get(UserUUID).(string)

			log := logger.Info().
				Str("remote_ip", vals.RemoteIP).
				Str("protocol", vals.Protocol).
				Str("host", vals.Host).
				Str("method", vals.Method).
				Str("uri", c.Request().RequestURI).
				Str("path", vals.RoutePath).
				Str("user_agent", vals.UserAgent).
				Int("status", vals.Status).
				Str("content_length", vals.ContentLength).
				Int64("response_size", vals.ResponseSize).
				Str("latency", vals.Latency.String())

			if vals.RequestID != "" {
				log.Str("request_id", vals.RequestID)
			}

			if userUUID != "" {
				log.Str(UserUUID, userUUID)
			}

			if vals.Referer != "" {
				log.Str("referer", vals.Referer)
			}

			if vals.Error != nil {
				log.Err(vals.Error)
			}

			log.Send()

			return nil
		},
	})
}
