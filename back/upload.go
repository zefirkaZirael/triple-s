package back

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
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
}
