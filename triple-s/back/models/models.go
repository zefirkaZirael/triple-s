package models

import "encoding/xml"

type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

// Buckets list for XML response
type ListBucketResponse struct {
	XMLName xml.Name `xml:"ListBuckets"`
	Buckets []Bucket `xml:"Bucket"`
}
