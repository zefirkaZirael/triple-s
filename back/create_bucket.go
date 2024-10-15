package back

import (
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"
)

var buckets = make(map[string]bool)

// Function to create buckets
func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		bucketName := r.URL.Path[1:]

		// bucket name valid check
		if bucketName == "" {
			http.Error(w, "Bucket name is required \n", http.StatusBadRequest)
			return
		}

		// Validate bucket name format
		if !isValidBucketName(bucketName) {
			http.Error(w, "Bucket name must be between 3-63 characters and can only contain lowercase letters, numbers, hyphens, and periods. \nMust not begin or end with a hyphen and must not contain two consecutive periods or dashes.\nMust not be formatted as an IP address (e.g., 192.168.0.1).\n", http.StatusBadRequest)
			return
		}

		// check if unique
		if buckets[bucketName] {
			http.Error(w, "Bucket already exists\n", http.StatusConflict)
			return
		}

		// Create bucket
		buckets[bucketName] = true
		if err := appendBucketToCSV(bucketName); err != nil {
			http.Error(w, "Failed to save bucket metadata\n", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Set status code to 200 OK
		fmt.Fprintf(w, "Bucket '%s' created successfully!\n", bucketName)
		return

	}
}

// Function to validate bucket name
func isValidBucketName(name string) bool {
	// Regular expression for bucket name validation
	validNamePattern := `^[a-z0-9](?:[a-z0-9.-]{1,61}[a-z0-9])?$`
	matched, _ := regexp.MatchString(validNamePattern, name)

	// Check for IP address format (simple check)
	if net.ParseIP(name) != nil {
		return false
	}

	return matched
}

func appendBucketToCSV(bucketName string) error {
	file, err := os.OpenFile("buckets.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{bucketName, time.Now().Format(time.RFC3339)}); err != nil {
		return err
	}

	return nil
}
