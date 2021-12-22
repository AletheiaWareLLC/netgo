package handler

import (
	"net/http"
)

func CacheControl(h http.Handler, cc string) http.Handler {
	return Header(h, "Cache-Control", cc)
}
