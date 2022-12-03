package main

import (
	"errors"
	"fmt"
	"net/http"
)

func main() {
    mux := http.NewServeMux()
    var flag bool
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("server: %s /%s\n", r.Method, r.URL.Path)
        if flag {
            http.Error(w, "I was Hacked!", http.StatusInternalServerError)
        }
        if r.URL.Path == "/hack" {
            flag = true
        }
    })
    server := http.Server{
        Addr:    fmt.Sprintf(":%d", 1080),
        Handler: mux,
    }
    if err := server.ListenAndServe(); err != nil {
        if !errors.Is(err, http.ErrServerClosed) {
            fmt.Printf("error running http server: %s\n", err)
        }
    }
}
