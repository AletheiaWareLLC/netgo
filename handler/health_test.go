package handler_test

import (
	"aletheiaware.com/netgo/handler"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	t.Run("Returns 200", func(t *testing.T) {
		mux := http.NewServeMux()
		handler.AttachHealthHandler(mux)
		request := httptest.NewRequest(http.MethodGet, "/health", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "", string(body))
	})
}
