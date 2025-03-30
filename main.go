package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hits: %d", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func main() {
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()	//Creates a server mux which routes http requests to handlers
	server := http.Server{	//Creates a server structure
		Handler: mux,
		Addr: ":8080",
	}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))	//Handles requests from /app/ endpoints, strips the /app and serves files in base directory
	mux.HandleFunc("/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("/reset", apiCfg.resetHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {	//Handles requests from /healthz endpoint
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
