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
	fmt.Println(objKey)
	if objKey == "objects.csv" {
		http.Error(w, "Cannot delete objects.csv", http.StatusForbidden)
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
		helpers.XMLResponse(w, http.StatusNotFound, "Object not found")
		return
	}

	if err := os.Remove(objectPath); err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to delete object")
		return
	}

	if err := removeObjectMetadata(w, bucketDir, objKey); err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to update metadata")
		return
	}
	if err := helpers.UpdateLastModified(bucketName); err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to update LastModified")
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
	// Set the permissions of the CSV file to read-only
	if err := os.Chmod(csvPath, 0444); err != nil { // Read-only for everyone
		return err
	}

	return nil
}
