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
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to read bucket metadata")
	}

	response := models.ListBucketResponse{Buckets: buckets}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	if err := xml.NewEncoder(w).Encode(response); err != nil {
		helpers.XMLResponse(w, http.StatusInternalServerError, "Failed to encode XML")
	}
}
