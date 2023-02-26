package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roman-wb/go-sample/pkg/test/helpers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
)

func Test_BodyLogger_WithoutAdditionalMiddlwares(t *testing.T) {
	logger, captureLog := helpers.NewLogger()

	e := echo.New()
	e.Use(BodyLogger(logger))
	e.GET("/", func(c echo.Context) error {
		return nil
	})

	server := httptest.NewServer(e)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.NotContains(t, captureLog.String(), "request_id")
	require.NotContains(t, captureLog.String(), "user_uuid")
	require.Contains(t, captureLog.String(), `"request_body":""`)
	require.Contains(t, captureLog.String(), `"response_body":""`)
}

func Test_BodyLogger_WithAdditionalMiddlwares(t *testing.T) {
	logger, captureLog := helpers.NewLogger()

	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(ExtractUserUUID())
	e.Use(BodyLogger(logger))
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "response body")
	})

	server := httptest.NewServer(e)
	defer server.Close()

	body := strings.NewReader(`request body`)
	req, err := http.NewRequest(http.MethodGet, server.URL, body)
	req.Header.Set(UserUUIDHeader, "test-uuid")
	require.NoError(t, err)
	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Regexp(t, `"request_id":".*"`, captureLog.String())
	require.Contains(t, captureLog.String(), `"user_uuid":"test-uuid"`)
	require.Contains(t, captureLog.String(), `"request_body":"request body"`)
	require.Contains(t, captureLog.String(), `"response_body":"response body"`)
}
