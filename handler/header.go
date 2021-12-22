package handler

import (
	"net/http"
)

func Header(h http.Handler, key, value string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(key, value)
		h.ServeHTTP(w, r)
	})
}
