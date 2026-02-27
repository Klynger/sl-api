package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
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
	conf := config.New()
	validator := validatorUtil.New()

	var logLevel gormlogger.LogLevel
	if conf.DB.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, conf.DB.Host, conf.DB.Username, conf.DB.Password, conf.DB.DBName, conf.DB.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})

	if err != nil {
		log.Fatal("DB connection start failure")
	}

	authStore := sessions.NewCookieStore([]byte(conf.StoreSessions.AuthSecret))

	r := router.New(db, validator, authStore, conf.Server.Debug, conf.StoreSessions.AuthMaxAgeSecs)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Server.Port),
		Handler:      r,
		ReadTimeout:  conf.Server.TimeoutRead,
		WriteTimeout: conf.Server.TimeoutWrite,
		IdleTimeout:  conf.Server.TimeoutIdle,
	}

	log.Println("Starting server ", s.Addr)
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}
}
