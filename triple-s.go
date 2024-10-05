package main

import (
	"log"
	"net/http"
	"os"
)

// curl -X PUT http://localhost:8080/mybucket1
// curl -X GET http://localhost:8080/

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: go run main.go <port> <directory_path>")
	}

	port := os.Args[1]
	// dir := os.Args[2]

	http.HandleFunc("/", listBuckets) // Point to the listBuckets handler
	// Set up the handler for creating buckets
	http.HandleFunc("/{BucketName}", createBucketHandler)

	// Start the server on port 8080
	log.Println("Starting server on port...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
