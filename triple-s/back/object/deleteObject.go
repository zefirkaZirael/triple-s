package object

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"triple-s/back/helpers"
)

func DeleteObject(w http.ResponseWriter, r *http.Request) {
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
	// Step3
	objectPath := filepath.Join(bucketDir, objKey)

	if _, err := os.Stat(objectPath); os.IsNotExist(err) {
		http.Error(w, "Object not found\n", http.StatusNotFound)
		return
	}

	if err := os.Remove(objectPath); err != nil {
		http.Error(w, "Failed to delete object\n", http.StatusInternalServerError)
		return
	}

	if err := removeObjectMetadata(w, bucketDir, objKey); err != nil {
		http.Error(w, "Failed to update metadata\n", http.StatusInternalServerError)
		return
	}
	if err := helpers.UpdateLastModified(bucketName); err != nil {
		http.Error(w, "Failed to update LastModified\n", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func removeObjectMetadata(w http.ResponseWriter, bucketDir, objKey string) error {
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
	defer tempFile.Close()

	reader := csv.NewReader(file)
	writer := csv.NewWriter(tempFile)

	// Copy all entries except the one to be removed
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return err
		}

		if record[0] != objKey { // Keep entries that don't match objKey
			if err := writer.Write(record); err != nil {
				return err
			}
		}
	}
	defer writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	// Replace the original CSV with the updated one
	if err := os.Rename(tempPath, csvPath); err != nil {
		return err
	}

	return nil
}
