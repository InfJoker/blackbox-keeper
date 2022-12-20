package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
	var flag bool
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("server: %s /%s\n", r.Method, r.URL.Path)
		if flag {
			http.Error(w, "I was Hacked!", http.StatusInternalServerError)
		}
		if r.URL.Path == "/hack" {
			flag = true
		}
	})
	if err := http.ListenAndServe(":1080", nil); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("error running http server: %s\n", err)
		}
	}
}
