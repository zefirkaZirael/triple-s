package bucket

import (
	"encoding/xml"
	"net/http"
	"triple-s/back/helpers"
	"triple-s/back/models"
)

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := helpers.ReadBucketMetadata()
	if err != nil {
		http.Error(w, "Failed to read bucket metadata\n", http.StatusInternalServerError)
		return
	}

	response := models.ListBucketResponse{Buckets: buckets}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	if err := xml.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode XML\n", http.StatusInternalServerError)
	}
}
