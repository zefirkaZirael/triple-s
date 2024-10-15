package back

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
)

// Bucket struct with metadata
type Bucket struct {
	Name         string `xml:"Name"`
	CreationTime string `xml:"CreationTime"`
}

// Buckets list for XML response
type ListBucketResponse struct {
	XMLName xml.Name `xml:"ListBuckets"`
	Buckets []Bucket `xml:"Bucket"`
}

// Read bucket metadata from a CSV file
func readBucketMetadata(filename string) ([]Bucket, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var buckets []Bucket
	for _, record := range records {
		name := record[0]

		creationTime := record[1]
		buckets = append(buckets, Bucket{Name: name, CreationTime: creationTime})
	}
	return buckets, nil
}

// Function to list all buckets
func ListBuckets(w http.ResponseWriter, r *http.Request) {
	// Make sure to have your CSV file named "buckets.csv"
	buckets, err := readBucketMetadata("buckets.csv")
	if err != nil {
		http.Error(w, "Failed to read bucket metadata\n", http.StatusInternalServerError)
		return
	}
	bucketsTrue := make(map[string]bool)
	for _, bucket := range buckets {
		if isValidBucketName(bucket.Name) {
			if bucketsTrue[bucket.Name] {
				http.Error(w, "Bucket already exists\n", http.StatusConflict)
				return
			}
			bucketsTrue[bucket.Name] = true
		} else {
			http.Error(w, fmt.Sprintf("Bucket name '%s' is invalid\n", bucket.Name), http.StatusBadRequest)
			return
		}
	}

	response := ListBucketResponse{Buckets: buckets}

	// Set the content type to XML
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	// For service
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	// Encode the response as XML and write it to the response writer
	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML\n", http.StatusInternalServerError)
	}
}
