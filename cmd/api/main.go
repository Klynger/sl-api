package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"sl-api/config"
)

func main() {
	c := config.New()

	mux := http.NewServeMux()
	mux.HandleFunc("/livez", livez)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      mux,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Println("Starting server ", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}

func livez(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello, world!")
}
