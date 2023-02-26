package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func Test_UserUUID_Empty(t *testing.T) {
	e := echo.New()
	e.Use(ExtractUserUUID())
	e.GET("/", func(c echo.Context) error {
		userUUID, _ := c.Get(UserUUID).(string)
		return c.HTML(http.StatusOK, userUUID)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.Empty(t, body)
}

func Test_UserUUID_Exists(t *testing.T) {
	e := echo.New()
	e.Use(ExtractUserUUID())
	e.GET("/", func(c echo.Context) error {
		userUUID, _ := c.Get(UserUUID).(string)
		return c.HTML(http.StatusOK, userUUID)
	})

	server := httptest.NewServer(e)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	req.Header.Set(UserUUIDHeader, "test-uuid")
	require.NoError(t, err)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	require.EqualValues(t, "test-uuid", string(body))
}
