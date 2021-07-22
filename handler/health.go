package handler

import (
	"net/http"
)

func AttachHealthHandler(m *http.ServeMux) {
	m.Handle("/health", Log(Health()))
}

func Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
