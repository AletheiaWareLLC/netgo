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

const CC = "max-age=60"

func TestStatic_File(t *testing.T) {
	mux := http.NewServeMux()
	fs := fstest.MapFS{
		"exists": {
			Data: []byte("hello, world"),
		},
	}
	handler.AttachStaticFSHandler(mux, fs, false, CC)
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

func TestStatic_Directory(t *testing.T) {
	t.Run("Returns 301 When Missing Trailing Slash", func(t *testing.T) {
		mux := http.NewServeMux()
		fs := fstest.MapFS{}
		handler.AttachStaticFSHandler(mux, fs, true, CC)
		request := httptest.NewRequest(http.MethodGet, "/static", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusMovedPermanently, result.StatusCode)
		assert.Equal(t, "<a href=\"/static/\">Moved Permanently</a>.\n\n", string(body))
	})
	t.Run("Returns 200 When Listable", func(t *testing.T) {
		mux := http.NewServeMux()
		fs := fstest.MapFS{
			"exists": {
				Data: []byte("hello, world"),
			},
		}
		handler.AttachStaticFSHandler(mux, fs, true, CC)
		request := httptest.NewRequest(http.MethodGet, "/static/", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "<pre>\n<a href=\"exists\">exists</a>\n</pre>\n", string(body))
	})
	t.Run("Returns 200 When Not Listable But index.html Exists", func(t *testing.T) {
		mux := http.NewServeMux()
		fs := fstest.MapFS{
			"index.html": {
				Data: []byte("hello, world"),
			},
		}
		handler.AttachStaticFSHandler(mux, fs, false, CC)
		request := httptest.NewRequest(http.MethodGet, "/static/", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "hello, world", string(body))
	})
	t.Run("Returns 404 When Not Listable And index.html Does Not Exist", func(t *testing.T) {
		mux := http.NewServeMux()
		fs := fstest.MapFS{
			"exists": {
				Data: []byte("hello, world"),
			},
		}
		handler.AttachStaticFSHandler(mux, fs, false, CC)
		request := httptest.NewRequest(http.MethodGet, "/static/", nil)
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, request)
		result := response.Result()
		body, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, result.StatusCode)
		assert.Equal(t, "404 page not found\n", string(body))
	})
}
