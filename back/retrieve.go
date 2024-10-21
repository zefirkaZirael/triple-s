package back

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func RetrieveObject(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path[1:], "/", 2)
	if len(parts) < 2 {
		http.Error(w, "Invalid request path. Format: /{BucketName}/{ObjectKey}\n", http.StatusBadRequest)
		return
	}

	bucketName := parts[0]
	objKey := parts[1]

	bucketDir := filepath.Join("data", bucketName)
	if _, err := os.Stat(bucketDir); os.IsNotExist(err) {
		http.Error(w, "Bucket not found\n", http.StatusNotFound)
		return
	}

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
}
