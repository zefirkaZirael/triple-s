package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"triple-s/back/bucket"
	"triple-s/back/object"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	// dir := falg.String("dir", ".", "Directory path")
	help := flag.Bool("help", false, "Show helper screen")
	flag.Usage = func() {
		fmt.Println("Simple Storage Service.")
		fmt.Println("\nUsage:")
		fmt.Println("    triple-s [-port <N>] [-dir <S>]")
		fmt.Println("    triple-s --help")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *port < 1 || *port > 65535 {
		log.Fatal("Invalid port number. Please choose a port between 1 and 65535.")
	}
	log.Println("Starting server on port...", *port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bucket.ListBuckets(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
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
		bucket.CreateBucketHandler(w, r)
	case http.MethodDelete:
		bucket.DeleteBucketHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func objectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		object.UploadObject(w, r) // Handle upload
	case http.MethodGet:
		object.RetrieveObject(w, r) // Handle retrieval
	case http.MethodDelete:
		object.DeleteObject(w, r) // Handle deletion
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
