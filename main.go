package main

import (
	"go-web/grun"
	"net/http"
)

func main() {
	r := grun.New()
	r.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})
	r.Run("8080")
}