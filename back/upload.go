package back

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadObject(w http.ResponseWriter, r *http.Request) {
	parts := filepath.SplitList(r.URL.Path[1:])
	if len(parts) < 2 {
		http.Error(w, "Invalid request path. Format: /{BucketName}/{ObjectKey}\n", http.StatusBadRequest)
		return
	}

	bucketName := parts[0]
	objKey := parts[1]

	// Check if the bucket exists
	if _, err := os.Stat("data/" + bucketName); os.IsNotExist(err) {
		http.Error(w, "Bucket not found\n", http.StatusNotFound)
		return
	}
	// Create the bucket directory if it doesn't exist
	bucketDir := filepath.Join("data", bucketName)
	objectPath := filepath.Join(bucketDir, objKey)

	file, err := os.Create(objectPath)
	if err != nil {
		http.Error(w, "Failed to save object\n", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	// Copy the binary data from request body to the file
	if _, err := io.Copy(file, r.Body); err != nil {
		http.Error(w, "Failed to write object data\n", http.StatusInternalServerError)
		return
	}

	// Save object metadata to CSV
	if err := appendObjectMetadata(bucketDir, objKey); err != nil {
		http.Error(w, "Failed to update object metadata\n", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Object '%s' uploaded successfully to bucket '%s'\n", objKey, bucketName)
}

func appendObjectMetadata(bucketDir, objKey string) err {
	csvPath := filepath.Join(bucketDir, "object.csv")

	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the object key and current timestamp
	if err := writer.Write([]string{objKey, time.Now().Format(time.RFC3339)}); err != nil {
		return err
	}
	return nil
}
