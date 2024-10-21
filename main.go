package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"triple-s/back"
)

// curl -X PUT http://localhost:8080/mybucket1
// curl -X GET http://localhost:8080/

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	// dir := os.Args[2]
	flag.Parse()
	if *port < 1 || *port > 65535 {
		log.Fatal("Invalid port number. Please choose a port between 1 and 65535.")
	}

	log.Println("Starting server on port...", *port)

	http.HandleFunc("/", back.ListBuckets)          // List all buckets
	http.HandleFunc("/{BucketName}", bucketHandler) // Create a bucket
	http.HandleFunc("/{BucketName}/{ObjectKey}", objectHandler)

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func bucketHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		back.CreateBucketHandler(w, r)
	case http.MethodDelete:
		back.DeleteBucketHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func objectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		back.UploadObject(w, r) // Handle upload
	case http.MethodGet:
		back.RetrieveObject(w, r) // Handle retrieval
	case http.MethodDelete:
		back.DeleteObject(w, r) // Handle deletion
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
