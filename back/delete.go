package back

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Path[1:]
	buckets, err := readBucketMetadata("buckets.csv")
	if err != nil {
		http.Error(w, "Failed to read bucket metadata\n", http.StatusInternalServerError)
		return
	}
	var bExist bool
	var bIndex int

	for i, bucket := range buckets {
		if bucket.Name == bucketName {
			bExist = true
			bIndex = i
			break
		}
	}

	if !bExist {
		http.Error(w, fmt.Sprintf("Bucket '%s' not found\n", bucketName), http.StatusNotFound)
		return
	}

	if !isBucketEmpty(bucketName) {
		http.Error(w, fmt.Sprintf("Bucket '%s' not empty\n", bucketName), http.StatusConflict)
		return
	}

	/////////////////!!?!?!!?!?
	buckets = append(buckets[:bIndex], buckets[bIndex+1:]...)

	err = saveBucketMetadata("buckets.csv", buckets)
	if err != nil {
		http.Error(w, "Failed to save bucket metadata\n", http.StatusInternalServerError)
		return
	}

	bucketDir := filepath.Join("data", bucketName)
	if err := os.RemoveAll(bucketDir); err != nil {
		http.Error(w, "Failed to delete bucket directory\n", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Bucket '%s' deleted successfully!\n", bucketName)
	w.WriteHeader(http.StatusNoContent)
}

func saveBucketMetadata(filename string, buckets []Bucket) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, bucket := range buckets {
		record := []string{bucket.Name, bucket.CreationTime}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func isBucketEmpty(bucketName string) bool {
	csvPath := filepath.Join("data", bucketName, "objects.csv")
	file, err := os.Open(csvPath)
	if err != nil {
		return true
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}
	for _, record := range records {
		if record[0] != "" {
			// object found so not empty
			return false
		}
	}
	// Default empty
	return true
}
