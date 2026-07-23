package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nodus-protocol/backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyLimit_AtLimit(t *testing.T) {
	const limit = int64(100)
	body := strings.Repeat("a", int(limit))

	r := gin.New()
	r.Use(BodyLimit(limit))
	r.POST("/test", func(c *gin.Context) {
		raw, err := c.GetRawData()
		if err != nil {
			utils.InternalServerError(c, err.Error())
			return
		}
		utils.OK(c, "ok", string(raw))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp utils.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, body, resp.Data)
}

func TestBodyLimit_OverLimit(t *testing.T) {
	const limit = int64(100)
	body := strings.Repeat("a", int(limit)+1)

	var handlerCalled bool
	r := gin.New()
	r.Use(BodyLimit(limit))
	r.POST("/test", func(c *gin.Context) {
		handlerCalled = true
		utils.OK(c, "ok", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
	r.ServeHTTP(w, req)

	assert.False(t, handlerCalled, "handler must not be called when body exceeds limit")
	require.Equal(t, http.StatusRequestEntityTooLarge, w.Code)

	var resp utils.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
	require.NotNil(t, resp.Error)
	assert.Equal(t, "PAYLOAD_TOO_LARGE", resp.Error.Code)
	assert.Equal(t, "request body exceeds maximum allowed size", resp.Error.Message)
}
