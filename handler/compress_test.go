package handler_test

import (
	"aletheiaware.com/netgo/handler"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompress(t *testing.T) {
	testhandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Hello World!"))
	})
	t.Run("ContentNoCompression", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/", handler.Compress(testhandler))
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "", response.HeaderMap.Get("Content-Encoding"))
		assert.Equal(t, "", response.HeaderMap.Get("Content-Length"))
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, "Hello World!", string(body))
	})
	t.Run("ContentGzip", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/", handler.Compress(testhandler))
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Accept-Encoding", "gzip")
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "gzip", response.HeaderMap.Get("Content-Encoding"))
		assert.Equal(t, "", response.HeaderMap.Get("Content-Length"))
		r, err := gzip.NewReader(result.Body)
		assert.Nil(t, err)
		body, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, "Hello World!", string(body))
	})
	t.Run("NoContentGzip", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.Handle("/", handler.Compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})))
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("Accept-Encoding", "gzip")
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		assert.Equal(t, http.StatusNoContent, result.StatusCode)
		assert.Equal(t, "", response.HeaderMap.Get("Content-Encoding"))
		assert.Equal(t, "", response.HeaderMap.Get("Content-Length"))
		assert.Equal(t, 0, response.Body.Len())
	})
}
