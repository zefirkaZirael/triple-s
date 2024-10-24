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
