package handler

import (
	"log"
	"net/http"
)

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL.Path, r.Header)
		h.ServeHTTP(w, r)
	})
}