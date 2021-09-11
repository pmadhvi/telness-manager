package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pmadhvi/telness-manager/handlers"
	"github.com/pmadhvi/telness-manager/postgres"
	"github.com/pmadhvi/telness-manager/service"
	"github.com/sirupsen/logrus"
)

func main() {
	// setup the log
	var log = logrus.New()
	log.SetOutput(os.Stdout)

	// read .env file for env variables
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file", err)
	}
	// read the port from env variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Info("port env variable not set, so using default port 8080")
		port = "8080"
	}
	dbname := os.Getenv("DB")
	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASS")
	dbhost := os.Getenv("DBHOST")
	dbport := os.Getenv("DBPORT")

	//Open db connection
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Error(errors.New("error connecting database, "), err)
	}
	var (
		subscriptionRepo = postgres.NewSubscriptionRepo(*log, db)
		subsvc           = service.SubscriptionSvc{SubscriptionRepo: subscriptionRepo}
	)

	// setup server and routes
	server := handlers.Server{
		Log:                 log,
		Port:                port,
		SubscriptionService: subsvc,
	}

	errorChan := make(chan error)

	// setup all the routes and start the server
	go func() {
		errorChan <- server.Start()
	}()

	// catch the exit signals and pass on errorChan
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		errorChan <- fmt.Errorf("Got quit signal: %s", <-quit)
	}()

	// get errors from chan and exit the application
	if err := <-errorChan; err != nil {
		log.Error(err)
		os.Exit(1)
	}

}
