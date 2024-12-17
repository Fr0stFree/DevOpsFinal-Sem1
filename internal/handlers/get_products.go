package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"project_sem/internal/config"
	"project_sem/internal/db"
	"project_sem/internal/serializers"
	"project_sem/internal/archivers"
)

func GetProducts(repo db.Repositorier) http.HandlerFunc {
	const errorResponseBody = "failed to load prices"
	const successContentType = "application/zip"
	const sucessContentDisposition = "attachment; filename=data.zip"
	const csvFileName = "data.csv"

	return func(w http.ResponseWriter, r *http.Request) {
		params, err := buildProductsFilter(r)
		if err != nil {
			log.Printf("failed to build filter params: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusBadRequest)
			return
		}
		products, err := repo.GetProducts(params)
		if err != nil {
			log.Printf("failed to load prices: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		serializedPrices, err := serializers.SerializeProduct(products)
		if err != nil {
			log.Printf("failed to serialize prices: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		archiver := archivers.New(archivers.ZipFmt)
		archive, err := archiver.Archive(w, csvFileName)
		if err != nil {
			log.Printf("failed to archive prices: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		defer archive.Close()
		_, err = archive.Write(serializedPrices.Bytes())
		if err != nil {
			log.Printf("failed to write prices to archive: %v\n", err)
			http.Error(w, errorResponseBody, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", successContentType)
		w.Header().Set("Content-Disposition", sucessContentDisposition)
	}
}

func buildProductsFilter(r *http.Request) (db.ProductsFilter, error) {
	params := db.ProductsFilter{}

	minCreateDate := r.URL.Query().Get("start")
	if minCreateDate != "" {
		minDate, err := time.Parse(config.DATE_FORMAT, minCreateDate)
		if err != nil {
			return params, err
		}
		params.MinCreateDate = minDate
	}

	maxCreateDate := r.URL.Query().Get("end")
	if maxCreateDate != "" {
		maxDate, err := time.Parse(config.DATE_FORMAT, maxCreateDate)
		if err != nil {
			return params, err
		}
		params.MaxCreateDate = maxDate
	}

	minPrice := r.URL.Query().Get("min")
	if minPrice != "" {
		price, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			return params, err
		}
		params.MinPrice = price
	}

	maxPrice := r.URL.Query().Get("max")
	if maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			return params, err
		}
		params.MaxPrice = price
	}

	return params, nil
}
