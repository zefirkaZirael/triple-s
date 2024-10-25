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
		helpers.XMLResponse(w, http.StatusBadRequest, "Bucket name is required")
		return
	}
	if !helpers.IsValidBucketName(bucketName) {
		helpers.XMLResponse(w, http.StatusBadRequest, "Bucket name must be between 3-63 characters and can only contain lowercase letters, numbers, hyphens, and periods.\nMust not begin or end with a hyphen and must not contain two consecutive periods or dashes.\nMust not be formatted as an IP address (e.g., 192.168.0.1).")
		return
	}
	for _, bucket := range buckets {
		if bucket.Name == bucketName {
			helpers.XMLResponse(w, http.StatusConflict, "Bucket already exists")
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
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to save bucket metadata")
		return
	}

	bucketDir := filepath.Join(helpers.BucketPath, bucketName)
	// Create bucket
	err := os.MkdirAll(bucketDir, 0o755)
	if err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to create bucket folder")
		return
	}

	w.WriteHeader(http.StatusOK) // Set status code to 200 OK
	helpers.XMLResponse(w, http.StatusOK, fmt.Sprintf("Bucket '%s' created successfully!", bucketName))
}
