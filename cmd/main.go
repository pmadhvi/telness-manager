package main

import (
	"database/sql"
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
		fmt.Errorf("Error loading .env file %v", err)
	}
	// read the port from env variable
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("port env variable not set, so using default port 8080")
		port = "8080"
	}
	dbname := os.Getenv("POSTGRES_DB")
	dbuser := os.Getenv("POSTGRES_USER")
	dbpass := os.Getenv("POSTGRES_PASSWORD")
	dbhost := os.Getenv("POSTGRES_HOST")
	dbport := os.Getenv("POSTGRES_PORT")

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
	if err = db.Ping(); err != nil {
		log.Errorf("Postgres ping error : (%v)", err)
	}
	defer db.Close()

	var (
		subscriptionRepo = postgres.NewSubscriptionRepo(db, log)
		subsvc           = service.SubscriptionSvc{SubscriptionRepo: subscriptionRepo, Log: log}
	)

	// setup server and routes
	server := handlers.Server{Log: log, Port: port, SubscriptionService: subsvc}

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
		fmt.Println(err)
		os.Exit(1)
	}

}
