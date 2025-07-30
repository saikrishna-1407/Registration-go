package main

import (
	"finalreg/config"
	"finalreg/handlers"
	"finalreg/internal/store"
	"finalreg/pkg/env"
	"fmt"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

const serviceName = "finalreg"

func main() {

	fmt.Println("Starting application...")

	fmt.Println("loading .env file...")

	//this will load environment variables from a .env file into the application's environment
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	} else {
		fmt.Println(".env file loaded successfully.")
	}

	fmt.Println("Loading environment variables into config struct...")

	//this will load environment variables directly into a custom configuration struct(conf of type config.Config)
	var conf config.Config
	if err := env.Load(&conf); err != nil {
		fmt.Println("Failed to load environment variables:", err)
		panic("Failed to load environment variables: " + err.Error())
	} else {
		fmt.Println("Environment variables loaded successfully.")
	}

	fmt.Println("Trimming DatabaseURI...")
	conf.DatabaseURI = strings.Trim(conf.DatabaseURI, "'")
	fmt.Printf("DatabaseURI after trimming: %s\n", conf.DatabaseURI)

	startService(&conf)
}

func startService(conf *config.Config) {
	fmt.Println("Connecting to PostgreSQL database...")

	psqlConn, err := connectPostgres(conf)
	if err != nil {
		fmt.Printf("Failed to connect to PostgreSQL. DatabaseURI: %s, Error: %v\n", conf.DatabaseURI, err)
		return
	}
	fmt.Println("Successfully connected to PostgreSQL database.")

	fmt.Println("Initializing Postgres store...")
	postgresStore := store.NewPostgresStore(psqlConn)
	fmt.Println("Postgres store initialized successfully.")

	fmt.Println("Creating service instance...")
	srv := &handlers.Service{
		ServiceName: serviceName,
		Config:      conf,
		Db:          postgresStore,
	}
	fmt.Println("Service instance created successfully.")

	fmt.Println("Setting up router...")
	router, err := handlers.SetupRouter(srv)
	if err != nil {
		fmt.Println(" Failed to setup router:", err)
	}
	fmt.Println(" Router setup complete")

	handlers.Handler()
	fmt.Println("Router setup successfully.")

	fmt.Println("Starting HTTP server on port 8000...")
	fmt.Println(http.ListenAndServe(":8000", router))
	fmt.Println("Server started successfully.")
}
