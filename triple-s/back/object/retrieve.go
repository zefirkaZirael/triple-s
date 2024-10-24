package object

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/back/helpers"
)

func RetrieveObject(w http.ResponseWriter, r *http.Request) {
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

	file, err := os.Open(objectPath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found\n", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Fail to open file\n", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	contentType := mime.TypeByExtension(filepath.Ext(objKey))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to send object data\n", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "Object '%s' retrieveed successfully from bucket '%s'\n", objKey, bucketName)
}
