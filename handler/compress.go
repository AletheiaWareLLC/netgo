package handler

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
	status int
}

func (w *gzipResponseWriter) Flush() {
	w.writer.Flush()
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.writer.Write(b)
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.status = status
	switch status {
	case http.StatusNoContent:
	case http.StatusNotModified:
	default:
		w.Header().Del("Content-Length")
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(status)
}

func Compress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		grw := &gzipResponseWriter{
			ResponseWriter: w,
			writer:         gzip.NewWriter(w),
		}
		defer func() {
			switch grw.status {
			case 0:
			case http.StatusNoContent:
			case http.StatusNotModified:
			default:
				if err := grw.writer.Close(); err != nil {
					log.Println(err)
				}
			}
		}()
		h.ServeHTTP(grw, r)
	})
}
