package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/back/bucket"
	"triple-s/back/helpers"
	"triple-s/back/object"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	dir := flag.String("dir", "data/", "Directory path")
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

	if err := os.MkdirAll(*dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}

	helpers.BucketPath = *dir
	helpers.BucketMetadataPath = filepath.Join(*dir, "buckets.csv")

	if _, err := os.Stat(helpers.BucketMetadataPath); os.IsNotExist(err) {
		file, err := os.Create(helpers.BucketMetadataPath)
		if err != nil {
			log.Fatalf("Failed to create buckets.csv: %v", err)
		}
		defer file.Close()
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bucket.ListBuckets(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/{BucketName}", bucketHandler) // Create a bucket
	http.HandleFunc("/{BucketName}/{ObjectKey}", objectHandler)

	log.Println("Starting server on port...", *port)
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
