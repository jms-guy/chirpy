package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"database/sql"
	"github.com/jms-guy/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" //Postgres driver, imported for side effects needed
)

func main() {
	err := godotenv.Load()	//Loads enviroment data from .env file
	if err != nil {
		fmt.Printf("Error loading .env file: %s", err)
		os.Exit(1)
	}
	dbURL := os.Getenv("DB_URL")	//Grabs database url
	db, err := sql.Open("postgres", dbURL)	//Opens database connection
	if err != nil {
		fmt.Printf("Error opening database connection: %s", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)	//Grabs generated sqlc queries

	apiCfg := &apiConfig{
		db: dbQueries,
	}

	mux := http.NewServeMux()	//Creates a server mux which routes http requests to handlers
	server := http.Server{	//Creates a server structure
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))	//Handles requests from /app/ endpoints, strips the /app and serves files in base directory
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)	//Handles server response to /admin/metrics	- displays visit count
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)	//Handles server response to /admin/reset - resets visit count
	mux.HandleFunc("POST /api/users", apiCfg.usersHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {	//Handles requests from /healthz endpoint
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")	//Sets response data
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	})

	fmt.Println("Listening...")
	err = server.ListenAndServe()	//Starts server
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}


