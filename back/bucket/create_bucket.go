package bucket

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"triple-s/back/helpers"
	"triple-s/back/models"
)

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[1:]
	buckets, _ := helpers.ReadBucketMetadata()
	if bucketName == "" {
		http.Error(w, "Bucket name is required \n", http.StatusBadRequest)
		return
	}
	if !helpers.IsValidBucketName(bucketName) {
		http.Error(w, "Bucket name must be between 3-63 characters and can only contain lowercase letters, numbers, hyphens, and periods. \nMust not begin or end with a hyphen and must not contain two consecutive periods or dashes.\nMust not be formatted as an IP address (e.g., 192.168.0.1).\n", http.StatusBadRequest)
		return
	}
	for _, bucket := range buckets {
		if bucket.Name == bucketName {
			http.Error(w, "Bucket already exists\n", http.StatusConflict)
			return
		}
	}
	creationTime := time.Now().Format(time.RFC3339)
	newBucket := models.Bucket{
		Name:             bucketName,
		CreationTime:     creationTime,
		LastModifiedTime: creationTime,
		Status:           "Active",
	}
	buckets = append(buckets, newBucket)
	if err := helpers.SaveBucketMetadata(buckets); err != nil {
		http.Error(w, "Failed to save bucket metadata\n", http.StatusInternalServerError)
		return
	}

	bucketDir := filepath.Join("data/", bucketName)
	// Create bucket
	err := os.MkdirAll(bucketDir, 0755)
	if err != nil {
		http.Error(w, "Failed to create bucket folder\n", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // Set status code to 200 OK
	fmt.Fprintf(w, "Bucket '%s' created successfully!\n", bucketName)
	return
}
