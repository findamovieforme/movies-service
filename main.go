package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/findamovieforme/movies-service/helpers"
	"github.com/findamovieforme/movies-service/routers"
)

func main() {
	helpers.InitEnv()
	router := routers.InitRouter()
	// Get the current GOMAXPROCS value
	currentProcs := runtime.GOMAXPROCS(0)
	fmt.Printf("Current GOMAXPROCS: %d\n", currentProcs)

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	go func() {
		log.Println("Server is running on port 8081...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("Server shutdown failed: %v", err)
	}
}
