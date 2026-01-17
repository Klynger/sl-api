package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"sl-api/api/router"
	"sl-api/config"
	validatorUtil "sl-api/util/validator"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

// @title			Shopping list API
// @version		1.0
// @description	An API to interact with the shopping list application
// @basePath		/v1
func main() {
	c := config.New()
	v := validatorUtil.New()

	var logLevel gormlogger.LogLevel
	if c.DB.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.Username, c.DB.Password, c.DB.DBName, c.DB.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})

	if err != nil {
		log.Fatal("DB connection start failure")
	}

	r := router.New(db, v)
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
