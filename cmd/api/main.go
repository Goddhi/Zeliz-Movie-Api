package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	
) 



///  application version number
const version = "1.0.0"

// struct for configuration settings for the application
type config struct {
	port int  // for the server
	env string // specifies the environment(dev, staging, production)
}

// an application struct to hold dependecncies for the http handlers, helpers, and middleware
type application struct {
	config config  // copy of the config strucy
	logger *log.Logger // 

}

func main() {
	// Declare an instance of the config struct
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	// the default value is set to 4000 and development
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse() /// allow us pass command line flags when running the application e.g go run ./cmd/api -port=33033 -env=production

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	// Declare an instance of the application struct, containing the config struct and
	// the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// setting up the router and handler
	// mux := http.NewServeMux()
	// mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//starting the http server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}



