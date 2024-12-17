package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"project_sem/internal/db"
	"project_sem/internal/serializers"
	"project_sem/internal/utils"
)

type PriceStats struct {
	TotalCount      int `json:"total_count"`
	DuplicateCount  int `json:"duplicates_count"`
	TotalItems      int `json:"total_items"`
	TotalCategories int `json:"total_categories"`
	TotalPrice      int `json:"total_price"`
}

func CreateProducts(repo *db.Repository) http.HandlerFunc {
	const errorResponseBody = "failed to upload prices"
	const successContentType = "application/json"

	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("failed to read incoming file: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		defer file.Close()

		formatType := r.URL.Query().Get("type")
		rc, err := unarchiveFile(file, formatType)
		if err != nil {
			log.Printf("failed to unarchive incoming file: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		stats := PriceStats{}
		products, totalCount, err := serializers.DeserializeProducts(rc)
		if err != nil {
			log.Printf("failed to parse prices from incoming file: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		stats.TotalCount = totalCount

		for _, product := range products {
			err = repo.CreateProduct(product)
			if err != nil {
				stats.DuplicateCount++
			} else {
				stats.TotalItems++
			}
		}
		totalPrice, totalCategories, err := repo.GetTotalPriceAndUniqueCategories()
		if err != nil {
			log.Printf("failed to get total price and unique categories: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		stats.TotalCategories = totalCategories
		stats.TotalPrice = int(totalPrice)

		w.Header().Set("Content-Type", successContentType)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(stats)
	}
}

func unarchiveFile(r io.Reader, fileType string) (io.ReadCloser, error) {
	switch fileType {
	case "zip":
		return utils.UnzipFile(r)
	case "tar":
		return utils.UntarFile(r)
	default:
		return utils.UnzipFile(r)
	}
}
