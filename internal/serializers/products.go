package serializers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"project_sem/internal/db"
)

func SerializeProduct(prices []db.Product) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)
	defer csvWriter.Flush()

	csvWriter.Write([]string{"id", "name", "category", "price", "create_date"})
	for _, price := range prices {
		record := []string{
			fmt.Sprintf("%d", price.ID),
			price.Name,
			price.Category,
			fmt.Sprintf("%.2f", price.Price),
			price.CreateDate.Format("2006-01-02"),
		}
		err := csvWriter.Write(record)
		if err != nil {
			return nil, err
		}
	}
	return &buffer, nil
}

func DeserializeProducts(r io.Reader) ([]db.Product, int, error) {
	prices := make([]db.Product, 0)
	csvReader := csv.NewReader(r)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, 0, err
	}
	totalItems := len(records) - 1
	for idx, record := range records {
		if idx == 0 {
			continue
		}
		price, err := validateProduct(record)
		if err != nil {
			log.Printf("price conversion failed %v\n", err)
			continue
		}
		prices = append(prices, price)
	}
	return prices, totalItems, nil
}

func validateProduct(record []string) (db.Product, error) {
	if len(record) != 5 {
		return db.Product{}, fmt.Errorf("invalid record %v", record)
	}
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return db.Product{}, fmt.Errorf("failed to convert id to int %v", err)
	}
	name := record[1]
	if name == "" {
		return db.Product{}, fmt.Errorf("name cannot be empty")
	}
	category := record[2]
	if category == "" {
		return db.Product{}, fmt.Errorf("category cannot be empty")
	}
	cost, err := strconv.ParseFloat(record[3], 64)
	if err != nil {
		return db.Product{}, fmt.Errorf("failed to convert cost to float %v", err)
	}
	create_date, err := time.Parse("2006-01-02", record[4])
	if err != nil {
		return db.Product{}, fmt.Errorf("failed to convert creation date %v", err)
	}
	price := db.Product{
		ID:         id,
		Name:       name,
		Category:   category,
		Price:      cost,
		CreateDate: create_date,
	}
	return price, nil
}