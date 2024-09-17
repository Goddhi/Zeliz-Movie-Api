package main

import (
	"context"
	"database/sql"
	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
	
	// Import the pq driver so that it can register itself with the database/sql
// package. Note that we alias this import to the blank identifier, to stop the Go
// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"  /// postgres driver
	

) 


///  application version number
const version = "1.0.0"

// struct for configuration settings for the application
type config struct {
	port int  // for the server
	env string // specifies the environment(dev, staging, production)
	db  struct {
		dsn	string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime string
	}
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

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSQL max open connection")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresSQL max idle connection")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresSQL max connection idle time")

	flag.Parse() /// allow us pass command line flags when running the application e.g go run ./cmd/api -port=33033 -env=production
	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	db, err := openDB(cfg) // creates the connection pool
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")
	
	
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
	err = srv.ListenAndServe()
	log.Fatal(err)
}


func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// the time.ParseDuration() function to convert the idle timeout duration string
// to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}

