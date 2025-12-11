package main

import (
	"fmt"
	"log"
	"net/http"

	"sl-api/api/router"
	"sl-api/config"
)

// @title			Shopping list API
// @version		1.0
// @description	An API to interact with the shopping list application
// @basePath		/v1
func main() {
	c := config.New()
	r := router.New()
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      r,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Println("Starting server ", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
