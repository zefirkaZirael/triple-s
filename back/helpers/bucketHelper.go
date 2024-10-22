package helpers

import (
	"encoding/csv"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"triple-s/back/models"
)

const bucketMetadataPath = "data/buckets.csv"

func ReadBucketMetadata() ([]models.Bucket, error) {
	file, err := os.Open(bucketMetadataPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var buckets []models.Bucket
	for _, record := range records {
		buckets = append(buckets, models.Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		})
	}
	return buckets, nil
}

func SaveBucketMetadata(buckets []models.Bucket) error {
	file, err := os.Create(bucketMetadataPath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, bucket := range buckets {
		record := []string{
			bucket.Name,
			bucket.CreationTime,
			bucket.LastModifiedTime,
			bucket.Status,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func IsValidBucketName(name string) bool {
	// Regular expression for bucket name validation
	validNamePattern := `^[a-z0-9](?:[a-z0-9.-]{1,61}[a-z0-9])?$`
	matched, _ := regexp.MatchString(validNamePattern, name)

	// Check for IP address format (simple check)
	if net.ParseIP(name) != nil {
		return false
	}

	return matched
}

func IsBucketEmpty(bucketName string) bool {
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

func UpdateLastModified(bucketName string) error {
	buckets, err := ReadBucketMetadata()
	if err != nil {
		return err
	}

	for i, bucket := range buckets {
		if bucket.Name == bucketName {
			buckets[i].LastModifiedTime = time.Now().Format(time.RFC3339)
			break
		}
	}

	return SaveBucketMetadata(buckets)
}
