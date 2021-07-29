package handler

import (
	"io/fs"
	"log"
	"net/http"
	"path"
)

func AttachStaticDirHandler(m *http.ServeMux, directory string) {
	m.Handle("/static/", Log(http.StripPrefix("/static/", StaticDir("html/static"))))
}

func AttachStaticFSHandler(m *http.ServeMux, fs fs.FS) {
	m.Handle("/static/", Log(http.StripPrefix("/static/", http.FileServer(http.FS(fs)))))
}

func StaticDir(directory string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "HEAD":
			fallthrough
		case "GET":
			http.ServeFile(w, r, path.Join(directory, r.URL.Path))
		default:
			log.Println("Unsupported method", r.Method)
		}
	})
}
