package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pmadhvi/telness-manager/client"
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
		log.Errorf("Error loading .env file %v", err)
	}
	// read the env variables from .env file
	port := os.Getenv("PORT")
	if port == "" {
		log.Info("port env variable not set, so using default port 8080")
		port = "8080"
	}
	dbname := os.Getenv("POSTGRES_DB")
	dbuser := os.Getenv("POSTGRES_USER")
	dbpass := os.Getenv("POSTGRES_PASSWORD")
	dbhost := os.Getenv("POSTGRES_HOST")
	dbport := os.Getenv("POSTGRES_PORT")
	ptsHost := os.Getenv("PTS_HOST")

	//Open db connection
	dbinfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbuser,
		dbpass,
		dbhost,
		dbport,
		dbname)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Errorf("error connecting database %v", err)
	}

	defer db.Close()

	var (
		subscriptionRepo = postgres.NewSubscriptionRepo(db, log)
		client           = client.NewClient(log, ptsHost)
		subsvc           = service.SubscriptionSvc{Log: log, SubscriptionRepo: subscriptionRepo, PtsClient: client}
	)

	// setup server and routes
	server := handlers.Server{Log: log, Port: port, SubscriptionService: subsvc}

	errorChan := make(chan error)
	quit := make(chan os.Signal, 1)
	// setup all the routes and start the server
	go func() {
		errorChan <- server.Start()
	}()

	// catch the exit signals and pass on errorChan
	go func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		errorChan <- fmt.Errorf("Got quit signal: %s", <-quit)
	}()

	// get errors from chan and exit the application
	if err := <-errorChan; err != nil {
		log.Info(err)
		os.Exit(1)
	}

}
