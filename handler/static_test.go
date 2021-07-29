package handler_test

import (
	"aletheiaware.com/netgo/handler"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestStatic(t *testing.T) {
	mux := http.NewServeMux()
	fs := fstest.MapFS{
		"exists": {
			Data: []byte("hello, world"),
		},
	}
	handler.AttachStaticFSHandler(mux, fs)
	t.Run("Returns 200 When File Exists", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/static/exists", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "hello, world", string(body))
	})
	t.Run("Returns 404 When File Does Not Exist", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/static/does-not-exist", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, result.StatusCode)
		assert.Equal(t, "404 page not found\n", string(body))
	})
}
