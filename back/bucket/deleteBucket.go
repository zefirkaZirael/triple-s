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
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to read bucket metadata")
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
		helpers.XMLResponse(w, http.StatusNotFound, fmt.Sprintf("Bucket '%s' not found", bucketName))
		return
	}

	if !helpers.IsBucketEmpty(bucketName) {
		buckets[bIndex].Status = "marked for deletion"
		if err := helpers.SaveBucketMetadata(buckets); err != nil {
			helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to update bucket status")
			return
		}
		if err := helpers.UpdateLastModified(bucketName); err != nil {
			helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to update LastModified")
			return
		}
		helpers.XMLResponse(w, http.StatusConflict, fmt.Sprintf("Bucket '%s' not empty", bucketName))
		return
	}

	buckets = append(buckets[:bIndex], buckets[bIndex+1:]...)

	err = helpers.SaveBucketMetadata(buckets)
	if err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to save bucket metadata")
		return
	}

	bucketDir := filepath.Join(helpers.BucketPath, bucketName)
	if err := os.RemoveAll(bucketDir); err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to delete bucket directory")
		return
	}
	// fmt.Fprintf(w, "Bucket '%s' deleted successfully!\n", bucketName)
	w.WriteHeader(http.StatusNoContent)
}
