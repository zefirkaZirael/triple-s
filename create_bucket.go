package main

import (
	"fmt"
	"net/http"
	"regexp"
)

// Function to create buckets
func createBucketHandler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "Bucket '%s' created successfully!\n", bucketName)
		return

	}
}

// Function to validate bucket name
func isValidBucketName(name string) bool {
	// Regular expression for bucket name validation
	validNamePattern := `^(?!.*--)(?!.*\.\.)(?!-)(?!.*-$)[a-z0-9.-]{3,63}$`
	matched, _ := regexp.MatchString(validNamePattern, name)

	// Check if the name is formatted as an IP address (simple check)
	ipPattern := `^\d{1,3}(\.\d{1,3}){3}$`
	ipMatched, _ := regexp.MatchString(ipPattern, name)

	fmt.Println(matched)
	fmt.Println(ipMatched)

	return matched && !ipMatched
}
