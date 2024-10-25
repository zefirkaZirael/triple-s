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
		helpers.XMLResponse(w, http.StatusNotFound, "File not found")
		return
	} else if err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Fail to open file")
		return
	}
	defer file.Close()

	contentType := mime.TypeByExtension(filepath.Ext(objKey))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	if _, err := io.Copy(w, file); err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to send object data")
		return
	}
	helpers.XMLResponse(w, http.StatusOK, fmt.Sprintf("Object '%s' retrieved successfully from bucket '%s'", objKey, bucketName))
}
