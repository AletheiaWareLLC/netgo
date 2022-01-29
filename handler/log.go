package handler

import (
	"aletheiaware.com/netgo"
	"net/http"
)

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		netgo.LogRequest(r)
		h.ServeHTTP(w, r)
	})
}
