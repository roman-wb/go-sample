package http

import (
	"net/http"
	"time"

	"github.com/roman-wb/go-sample/pkg/config"
	appMiddleware "github.com/roman-wb/go-sample/pkg/servers/http/middleware"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func New(logger *zerolog.Logger, cfg *config.Config) *http.Server {
	handler := echo.New()

	handler.Use(echoMiddleware.Recover())
	handler.Use(echoMiddleware.RequestID())
	handler.Use(appMiddleware.ExtractUserUUID())
	handler.Use(appMiddleware.Logger(logger))
	handler.Use(appMiddleware.BodyLogger(logger))

	handler.GET("/", func(c echo.Context) error {
		time.Sleep(5 * time.Second) //nolint:gomnd

		return c.String(http.StatusOK, "Hello, World!")
	})

	return &http.Server{
		Addr:         cfg.Server,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
