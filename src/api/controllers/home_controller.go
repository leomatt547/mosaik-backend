package controllers

import (
	"fmt"
	"net/http"

	_ "github.com/heroku/x/hmetrics/onload"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	p := "." + r.URL.Path
	if p == "./" {
		p = "./index.html"
	}
	http.ServeFile(w, r, p)
}
