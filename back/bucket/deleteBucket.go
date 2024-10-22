package bucket

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/back/helpers"
)

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[1:]
	buckets, err := helpers.ReadBucketMetadata()
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

	if !helpers.IsBucketEmpty(bucketName) {
		buckets[bIndex].Status = "marked for deletion"
		if err := helpers.SaveBucketMetadata(buckets); err != nil {
			http.Error(w, "Failed to update bucket status\n", http.StatusInternalServerError)
			return
		}
		http.Error(w, fmt.Sprintf("Bucket '%s' not empty\n", bucketName), http.StatusConflict)
		return
	}

	buckets = append(buckets[:bIndex], buckets[bIndex+1:]...)

	err = helpers.SaveBucketMetadata(buckets)
	if err != nil {
		http.Error(w, "Failed to save bucket metadata\n", http.StatusInternalServerError)
		return
	}

	bucketDir := filepath.Join("data", bucketName)
	if err := os.RemoveAll(bucketDir); err != nil {
		http.Error(w, "Failed to delete bucket directory\n", http.StatusInternalServerError)
		return
	}
	// fmt.Fprintf(w, "Bucket '%s' deleted successfully!\n", bucketName)
	w.WriteHeader(http.StatusNoContent)
}
