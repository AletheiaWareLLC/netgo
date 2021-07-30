package handler

import (
	"io/fs"
	"net/http"
	"path/filepath"
)

func AttachStaticDirHandler(m *http.ServeMux, directory string, listable bool) {
	AttachStaticHTTPFSHandler(m, http.Dir(directory), listable)
}

func AttachStaticFSHandler(m *http.ServeMux, fs fs.FS, listable bool) {
	AttachStaticHTTPFSHandler(m, http.FS(fs), listable)
}

func AttachStaticHTTPFSHandler(m *http.ServeMux, fs http.FileSystem, listable bool) {
	m.Handle("/static/", Log(http.StripPrefix("/static/", StaticFS(fs, listable))))
}

func StaticDir(directory string, listable bool) http.Handler {
	return StaticFS(http.Dir(directory), listable)
}

func StaticFS(filesystem http.FileSystem, listable bool) http.Handler {
	return http.FileServer(&staticFS{filesystem, listable})
}

type staticFS struct {
	fs       http.FileSystem
	listable bool
}

func (s *staticFS) Open(path string) (http.File, error) {
	file, err := s.fs.Open(path)
	if err != nil {
		return nil, err
	}
	if !s.listable {
		// Check if path is a directory
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			// Check if index.html exists
			index, err := s.fs.Open(filepath.Join(path, "index.html"))
			if err != nil {
				// Close directory
				if err := file.Close(); err != nil {
					return nil, err
				}
				return nil, err
			}
			// Close index
			if err := index.Close(); err != nil {
				return nil, err
			}
		}
	}

	return file, nil
}
