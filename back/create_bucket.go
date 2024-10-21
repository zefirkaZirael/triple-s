package back

import (
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	//r.PathValue("Bucke")
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method\n", http.StatusMethodNotAllowed)
		return
	}
	bucketName := r.URL.Path[1:]

	if bucketName == "" {
		http.Error(w, "Bucket name is required \n", http.StatusBadRequest)
		return
	}

	if !isValidBucketName(bucketName) {
		http.Error(w, "Bucket name must be between 3-63 characters and can only contain lowercase letters, numbers, hyphens, and periods. \nMust not begin or end with a hyphen and must not contain two consecutive periods or dashes.\nMust not be formatted as an IP address (e.g., 192.168.0.1).\n", http.StatusBadRequest)
		return
	}

	bucketDir := filepath.Join("data/", bucketName)
	if _, err := os.Stat(bucketDir); !os.IsNotExist(err) {
		http.Error(w, "Bucket already exists\n", http.StatusConflict)
		return
	}

	// Create bucket
	err := os.MkdirAll(bucketDir, 0755)
	if err != nil{
		http.Error(w, "Failed to create bucket folder\n", http.StatusInternalServerError)
		return
	}
	
	if err := appendBucketToCSV(bucketName); err != nil {
		http.Error(w, "Failed to save bucket metadata\n", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK) // Set status code to 200 OK
	fmt.Fprintf(w, "Bucket '%s' created successfully!\n", bucketName)
	return
}

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

	err1 := writer.Write([]string{bucketName, time.Now().Format(time.RFC3339)})
	return err1
}
