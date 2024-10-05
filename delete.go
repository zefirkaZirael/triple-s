package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
)

func deleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[1:]

	buckets, err := readBucketMetadata("buckets.csv")
	if err != nil {
		http.Error(w, "Failed to read bucket metadata\n", http.StatusInternalServerError)
		return
	}
	var bExist bool
	var bIndex int

	for i, bucket := range buckets {
		if bucket.Name == bucketName {
			bExist = true
			bIndex = i
			break
		}
	}

	if !bExist {
		http.Error(w, fmt.Sprintf("Bucket '%s' not found\n", bucketName), http.StatusNotFound)
		return
	}

	if !isBucketEmpty(bucketName) {
		http.Error(w, fmt.Sprintf("Bucket '%s' not found\n", bucketName), http.StatusConflict)
		return
	}
	if bIndex == len(buckets)-1 {
		buckets = append(buckets[:bIndex])
	} else {
		buckets = append(buckets[:bIndex], buckets[bIndex+1]...)
	}
	err = saveBucketMetadata("buckets.csv", buckets)
	if err != nil {
		http.Error(w, "Failed to save bucket metadata\n", http.StatusInternalServerError)
		return
	}

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func saveBucketMetadata(filename string, buckets []Bucket) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	defer writer.Flush()

	for _, bucket := range buckets {
		record := []string{bucket.Name, bucket.CreationTime}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func isBucketEmpty(bucketNmae string) bool {
	return true
}
