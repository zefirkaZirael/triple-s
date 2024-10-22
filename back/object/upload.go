package object

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"triple-s/back/helpers"
)

func UploadObject(w http.ResponseWriter, r *http.Request) {
	// Step 1
	bucketName, objKey, err := helpers.ParseBucketAndObjectKey(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Step 2
	bucketDir, err := helpers.BucketExists(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Step 3
	objectPath := filepath.Join(bucketDir, objKey)
	if _, err := os.Stat(objectPath); !os.IsNotExist(err) {
		http.Error(w, "Object already exists\n", http.StatusConflict)
		return
	}

	file, err := os.Create(objectPath)
	if err != nil {
		http.Error(w, "Failed to save object\n", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	size, err := io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Failed to write object data\n", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(make([]byte, 512))
	r.Body.Read(make([]byte, 512))
	// Save object metadata to CSV
	if err := appendObjectMetadata(bucketDir, objKey, size, contentType); err != nil {
		http.Error(w, "Failed to update object metadata\n", http.StatusInternalServerError)
		return
	}
	if err := helpers.UpdateLastModified(bucketName); err != nil {
		http.Error(w, "Failed to update LastModified\n", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Object '%s' uploaded successfully to bucket '%s'\n", objKey, bucketName)
}

func appendObjectMetadata(bucketDir, objKey string, size int64, contentType string) error {
	csvPath := filepath.Join(bucketDir, "objects.csv")

	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	record := []string{
		objKey,
		fmt.Sprintf("%d", size),
		contentType,
		time.Now().Format(time.RFC3339),
	}
	if err := writer.Write(record); err != nil {
		return err
	}
	return nil
}
