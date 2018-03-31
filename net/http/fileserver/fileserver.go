package fileserver

import (
	"net/http"
)

func isIndex(url string) bool {
	if url[len(url)-1] != '/' {
		return true
	}
	return false
}

func FileHandler(w http.ResponseWriter, r *http.Request) {
	if isIndex(r.URL.Path) {
		http.ServeFile(w, r, r.URL.Path[1:])
	} else {
		http.NotFound(w, r)
	}
}
