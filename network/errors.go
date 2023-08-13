package network

import (
	"fmt"
	"net/http"
)

func FourOhFour(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, s-maxage=604800")
	fmt.Fprint(w, "<h1>404</h1>");
	return
}