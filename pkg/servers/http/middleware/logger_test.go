package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/roman-wb/go-sample/pkg/test/helpers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"
)

func Test_Logger_WithoutAdditionalMiddlwares(t *testing.T) {
	logger, captureLog := helpers.NewLogger()

	e := echo.New()
	e.Use(Logger(logger))
	e.GET("/test", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "hello")
	})

	server := httptest.NewServer(e)
	defer server.Close()

	body := strings.NewReader(`request body`)
	req, err := http.NewRequest(http.MethodGet, server.URL+"/test?f=1", body)
	require.NoError(t, err)
	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Contains(t, captureLog.String(), `"level":"info"`)
	require.Contains(t, captureLog.String(), `"remote_ip":"127.0.0.1"`)
	require.Contains(t, captureLog.String(), `"protocol":"HTTP/1.1"`)
	require.Regexp(t, `"host":"127.0.0.1:.*"`, captureLog.String())
	require.Contains(t, captureLog.String(), `"method":"GET"`)
	require.Contains(t, captureLog.String(), `"uri":"/test?f=1"`)
	require.Contains(t, captureLog.String(), `"path":"/test"`)
	require.Contains(t, captureLog.String(), `"user_agent":"Go-http-client/1.1"`)
	require.Contains(t, captureLog.String(), `"status":200`)
	require.Contains(t, captureLog.String(), `"content_length":"12"`)
	require.Contains(t, captureLog.String(), `"response_size":5`)
	require.Regexp(t, `"latency":".*"`, captureLog.String())

	require.NotContains(t, captureLog.String(), `request_id`)
	require.NotContains(t, captureLog.String(), `referer`)
	require.NotContains(t, captureLog.String(), UserUUID)
	require.NotContains(t, captureLog.String(), `err`)
}

func Test_Logger_WithAdditionalMiddlwares(t *testing.T) {
	logger, captureLog := helpers.NewLogger()

	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(ExtractUserUUID())
	e.Use(Logger(logger))
	e.GET("/test", func(c echo.Context) error {
		return errors.New("some error")
	})

	server := httptest.NewServer(e)
	defer server.Close()

	body := strings.NewReader(`request body`)
	req, err := http.NewRequest(http.MethodGet, server.URL+"/test?f=1", body)
	req.Header.Set(UserUUIDHeader, "test-uuid")
	req.Header.Set("Referer", "some-referer")
	require.NoError(t, err)
	_, err = http.DefaultClient.Do(req)
	require.NoError(t, err)

	require.Contains(t, captureLog.String(), `"level":"info"`)
	require.Contains(t, captureLog.String(), `"remote_ip":"127.0.0.1"`)
	require.Contains(t, captureLog.String(), `"protocol":"HTTP/1.1"`)
	require.Regexp(t, `"host":"127.0.0.1:.*"`, captureLog.String())
	require.Contains(t, captureLog.String(), `"method":"GET"`)
	require.Contains(t, captureLog.String(), `"uri":"/test?f=1"`)
	require.Contains(t, captureLog.String(), `"path":"/test"`)
	require.Contains(t, captureLog.String(), `"user_agent":"Go-http-client/1.1"`)
	require.Contains(t, captureLog.String(), `"status":200`)
	require.Contains(t, captureLog.String(), `"content_length":"12"`)
	require.Contains(t, captureLog.String(), `"response_size":0`)
	require.Regexp(t, `"latency":".*"`, captureLog.String())

	require.Contains(t, captureLog.String(), `"user_uuid":"test-uuid"`)
	require.Contains(t, captureLog.String(), `"error":"some error"`)
	require.Regexp(t, `"request_id":".*"`, captureLog.String())
	require.Contains(t, captureLog.String(), `"referer":"some-referer"`)
}
