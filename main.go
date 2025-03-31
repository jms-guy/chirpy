package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()	//Creates a server mux which routes http requests to handlers
	server := http.Server{	//Creates a server structure
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))	//Handles requests from /app/ endpoints, strips the /app and serves files in base directory
	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)	//Handles server response to /admin/metrics	- displays visit count
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)	//Handles server response to /admin/reset - resets visit count
	mux.HandleFunc("POST /api/validate_chirp", validateHandler)
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {	//Handles requests from /healthz endpoint
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")	//Sets response data
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "OK")
	})

	fmt.Println("Listening...")
	err := server.ListenAndServe()	//Starts server
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}


