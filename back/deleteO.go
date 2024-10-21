package back

import (
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DeleteObject(w http.ResponseWriter, r *http.Request) {
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
	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		http.Error(w, "Object not found\n", http.StatusNotFound)
		return
	}

	if err := os.Remove(objectPath); err != nil {
		http.Error(w, "Failed to delete object\n", http.StatusInternalServerError)
		return
	}

	if err := removeObjectMetadata(bucketDir, objKey); err != nil {
		http.Error(w, "Failed to update metadata\n", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func removeObjectMetadata(bucketDir, objKey string) error {
	csvPath := filepath.Join(bucketDir, "objects.csv")
	tempPath := filepath.Join(bucketDir, "temp.csv")

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}

	defer file.Close()

	tempFile, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	// Copy all entries except the one to be removed
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if record[0] != objKey { // Keep entries that don't match objKey
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	// Replace the original CSV with the updated one
	if err := os.Rename(tempPath, csvPath); err != nil {
		return err
	}

	return nil
}
