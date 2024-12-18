package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"project_sem/internal/archivers"
	"project_sem/internal/db"
	"project_sem/internal/serializers"
)

type PriceStats struct {
	TotalCount      int `json:"total_count"`
	DuplicateCount  int `json:"duplicates_count"`
	TotalItems      int `json:"total_items"`
	TotalCategories int `json:"total_categories"`
	TotalPrice      int `json:"total_price"`
}

func CreateProducts(repo db.Repositorier) http.HandlerFunc {
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

		fmt := r.URL.Query().Get("type")
		archiver := archivers.New(archivers.Format(fmt))
		rc, err := archiver.Extract(file)
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

		transaction, err := repo.Begin()
		for _, product := range products {
			err = repo.CreateProduct(product)
			if err == nil {
				stats.TotalItems++
				continue
			}
			if db.IsDuplicateError(err) {
				stats.DuplicateCount++
				continue
			}
			transaction.Rollback()
			log.Printf("failed to save product: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		err = transaction.Commit()
		if err != nil {
			log.Printf("failed to commit transaction: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
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
