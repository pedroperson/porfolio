package handler

import (
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}
