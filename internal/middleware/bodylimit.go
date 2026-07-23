package middleware

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nodus-protocol/backend/internal/utils"
)

func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			var maxBytesErr *http.MaxBytesError
			if errors.As(err, &maxBytesErr) {
				utils.PayloadTooLarge(c)
				c.Abort()
				return
			}
			utils.InternalServerError(c, "failed to read request body")
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))

		c.Next()
	}
}
