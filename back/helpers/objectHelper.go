package helpers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ParseBucketAndObjectKey(r *http.Request) (string, string, error) {
	parts := strings.SplitN(r.URL.Path[1:], "/", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid request path. Format: /{BucketName}/{ObjectKey}\n")
	}
	return parts[0], parts[1], nil
}

func BucketExists(bucketName string) (string, error) {
	bucketDir := filepath.Join(BucketPath, bucketName)
	if _, err := os.Stat(bucketDir); os.IsNotExist(err) {
		return "", fmt.Errorf("bucket not found\n")
	}
	return bucketDir, nil
}

func XMLResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	fmt.Fprintf(w, "<response>\n\t<message>%s</message>\n</response>\n", message)
}
