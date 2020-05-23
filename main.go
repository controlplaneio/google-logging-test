// Sample stdlogging writes log.Logger logs to the Stackdriver Logging.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/logging"
)

var logger *log.Logger

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := "kubesim-dev-20200331"

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Sets the name of the log to write to.
	logName := "google-cloud-logging-appengine-tes"

	logger = client.Logger(logName).StandardLogger(logging.Info)

	// Logs "hello world", log entry is visible at
	// Stackdriver Logs.
	logger.Println("hello world")
	http.HandleFunc("/", handle)
	http.HandleFunc("/_ah/health", healthCheckHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello world!")
	logger.Print("Received request to /")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Print("Recieved request to healthcheck")
	fmt.Fprint(w, "ok")
}
