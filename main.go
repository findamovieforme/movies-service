package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/movierecuh/movies-service/helpers"
	"github.com/movierecuh/movies-service/routers"
)

func main() {
	helpers.InitEnv()
	router := routers.InitRouter()
	log.Println("Server is running on port 8081...")
	// Get the current GOMAXPROCS value
	currentProcs := runtime.GOMAXPROCS(0)
	fmt.Printf("Current GOMAXPROCS: %d\n", currentProcs)
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}
